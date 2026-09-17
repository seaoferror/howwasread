package service

import (
	"backend/common/payload"
	"backend/onlineconversation/internal/dto"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func (s *Service) CreateConversation(
	ctx context.Context,
	memberId uuid.UUID,
	novel,
	shortStory,
	poem,
	play,
	film,
	writtenBy,
	rule string,
	capacity int,
	t time.Time,
	length time.Duration,
) (map[string]uuid.UUID, error) {
	conversationId, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid v7 for online conversation id", "err", err)
		return nil, err
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		slog.Error("fail to start transaction for create conversation", "err", err)
		return nil, err
	}
	defer tx.Rollback()

	err = s.repository.InsertConversation(ctx, tx, conversationId, novel, shortStory, poem, play, film, writtenBy, rule, capacity, t, length, memberId)
	if err != nil {
		return nil, err
	}
	err = s.repository.InsertRegistrant(ctx, tx, conversationId, memberId)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("fail to commit transaction for create conversation", "err", err)
		return nil, err
	}
	slog.Info("success to create conversation")
	s.producer.PushMessage("search",
		nil,
		payload.Marshal(dto.OnlineConversationDocument{
			Id:         conversationId,
			Novel:      novel,
			ShortStory: shortStory,
			Poem:       poem,
			Play:       play,
			Film:       film,
			WrittenBy:  writtenBy,
			Time:       t,
		}),
		[]sarama.RecordHeader{
			{Key: []byte("type"), Value: []byte("onlineconversation")},
		},
	)
	return map[string]uuid.UUID{"conversationId": conversationId}, nil
}

func (s *Service) GetConversations(ctx context.Context, page int, t time.Time) ([]dto.OnlineConversationFeedResponse, error) {
	resp := []dto.OnlineConversationFeedResponse{}

	items, err := s.repository.FindConversations(ctx, s.repository.Tx(), page, t)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		resp = append(resp, dto.OnlineConversationFeedResponse{
			Id:         item.Id,
			Novel:      item.Novel,
			ShortStory: item.ShortStory,
			Poem:       item.Poem,
			Play:       item.Play,
			Film:       item.Film,
			WrittenBy:  item.WrittenBy,
			Time:       item.Time,
		})
	}
	slog.Info("success to get conversation", "resp", resp)
	return resp, nil
}
func (s *Service) PublishConversationSignal(fromId uuid.UUID, toIds [][]byte, signal []byte) error {
	value := payload.Marshal(payload.OnlineConversationSignal{
		FromId: fromId[:],
		ToIds:  toIds,
		Signal: signal,
	})
	s.producer.PushMessage("conversation-signal", nil, value, nil)
	return nil
}

func (s *Service) GetConversationDetail(ctx context.Context, conversationId, memberId uuid.UUID) (*dto.OnlineConversationDetailResponse, error) {
	c, err := s.repository.FindConversation(ctx, s.repository.Tx(), conversationId)
	if err != nil {
		return nil, err
	}
	var isRegistrant bool
	for _, r := range c.RegistrantIds {
		if r == memberId {
			isRegistrant = true
			break
		}
	}
	canEnter := true
	if time.Now().UTC().Before(c.Time.Add(-15 * time.Minute)) {
		canEnter = false
	}
	if time.Now().UTC().Before(c.Time.Add(10*time.Minute)) && !isRegistrant {
		canEnter = false
	}
	for _, b := range c.BanIds {
		if b == memberId {
			canEnter = false
			break
		}
	}
	var isNotificationScheduled bool
	for _, n := range c.NotificationIds {
		if n == memberId {
			isNotificationScheduled = true
		}
	}
	resp := dto.OnlineConversationDetailResponse{
		Id:                      c.Id,
		Novel:                   c.Novel,
		ShortStory:              c.ShortStory,
		Poem:                    c.Poem,
		Play:                    c.Play,
		Film:                    c.Film,
		WrittenBy:               c.WrittenBy,
		Rule:                    c.Rule,
		Capacity:                c.Capacity,
		Time:                    c.Time,
		Length:                  c.Length.String(),
		CanEnter:                canEnter,
		ModeratorIds:            c.ModeratorIds,
		IsRegistrant:            isRegistrant,
		IsNotificationScheduled: isNotificationScheduled,
	}
	return &resp, nil
}

func (s *Service) BanParticipant(ctx context.Context, modId, conversationId, banId uuid.UUID) error {
	mIds, err := s.repository.FindModeratorIds(ctx, s.repository.Tx(), conversationId)
	if err != nil {
		return err
	}
	isMod := false
	for _, mId := range mIds {
		if mId == modId {
			isMod = true
			break
		}
	}
	if !isMod {
		return errors.New("you cannot ban")
	}
	err = s.repository.AddBanId(ctx, s.repository.Tx(), conversationId, banId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ReportOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	ids, err := s.repository.FindReporterIds(ctx, s.repository.Tx(), conversationId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if id == memberId {
			return nil
		}
	}
	if len(ids) > 5 {
		err = s.repository.DeleteOnlineConversation(ctx, s.repository.Tx(), conversationId)
		if err != nil {
			return err
		}
	}
	err = s.repository.AddReporterId(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RegisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		slog.Error("fail to start transaction for register conversation", "err", err)
		return err
	}
	defer tx.Rollback()

	currentRegistrants, capacity, err := s.repository.LockConversationRegistrants(ctx, tx, conversationId)
	if err != nil {
		return err
	}
	if currentRegistrants >= capacity {
		return errors.New("already fully registered")
	}

	err = s.repository.InsertRegistrant(ctx, tx, conversationId, memberId)
	if err != nil {
		return err
	}

	err = s.repository.IncrementRegistrants(ctx, tx, conversationId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("fail to commit transaction for register conversation", "err", err)
		return err
	}
	return nil
}

func (s *Service) DeregisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		slog.Error("fail to start transaction for deregister conversation", "err", err)
		return err
	}
	defer tx.Rollback()

	err = s.repository.RemoveRegistrantId(ctx, tx, conversationId, memberId)
	if err != nil {
		return err
	}

	err = s.repository.DecrementRegistrants(ctx, tx, conversationId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("fail to commit transaction for deregister conversation", "err", err)
		return err
	}
	return nil
}

func (s *Service) GenerateTurn() *dto.GetTurnResponse {
	res := &dto.GetTurnResponse{
		Uris: []string{
			fmt.Sprintf("turn:%s:3478?transport=udp", s.turnRealm),
			fmt.Sprintf("turn:%s:5349?transport=tcp", s.turnRealm),
		},
		Username: fmt.Sprintf("%d", time.Now().Add(2*time.Hour).Unix()),
	}
	mac := hmac.New(sha1.New, []byte(s.turnSecret))
	mac.Write([]byte(res.Username))
	res.Credential = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return res
}

func (s *Service) ScheduleNotification(ctx context.Context, memberId, conversationId uuid.UUID) error {
	c, err := s.repository.FindConversation(ctx, s.repository.Tx(), conversationId)
	if err != nil {
		return err
	}
	err = s.repository.AddNotificationId(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return err
	}
	p := payload.NotificationScheduling{
		PartitionId: conversationId,
		KeyId:       memberId,
	}
	if len(c.NotificationIds) == 0 {
		aboutRaw := []rune(c.Novel + c.Play + c.Poem + c.ShortStory + c.Film + c.WrittenBy)
		if len(aboutRaw) > 6 {
			aboutRaw = []rune(string(aboutRaw[:6]) + "...")
		}
		p.ScheduledTime = c.Time.Add(-15 * time.Minute).UnixMilli()
		p.Contents = map[int]string{0: string(aboutRaw)}
		p.Type = "online-conversation"
	}
	s.producer.PushMessage("scheduled-notification", nil,
		payload.Marshal(p),
		nil)
	return nil
}

func (s *Service) CancelNotification(ctx context.Context, memberId, conversationId uuid.UUID) error {
	err := s.repository.RemoveNotificationId(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return err
	}
	s.producer.PushMessage("scheduled-notification", nil,
		payload.Marshal(payload.NotificationScheduling{
			PartitionId: conversationId,
			KeyId:       memberId,
			Type:        "cancel",
		}),
		nil,
	)
	return nil
}
