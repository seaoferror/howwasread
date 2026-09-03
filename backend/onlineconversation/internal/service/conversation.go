package service

import (
	"backend/common/payload"
	"backend/onlineconversation/internal/dto"
	"bytes"
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
	err = s.repository.SaveConversation(
		ctx,
		memberId,
		conversationId,
		novel,
		shortStory,
		poem,
		play,
		film,
		writtenBy,
		rule,
		capacity,
		t,
		length,
	)
	if err != nil {
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

	items, err := s.repository.FindConversations(ctx, page, t)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		resp = append(resp, dto.OnlineConversationFeedResponse{
			Id:         uuid.UUID(item.Id.Data),
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

func (s *Service) GetParticipantsWithoutMe(ctx context.Context, conversationId string, memberId uuid.UUID) ([]uuid.UUID, error) {
	pidRaws, err := s.repository.FindParticipantIds(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	pids := make([]uuid.UUID, 0, len(pidRaws))
	for _, pidRaw := range pidRaws {
		pid, err1 := uuid.FromBytes([]byte(pidRaw))
		if err1 != nil {
			slog.Error("fail to parse uuid from pidRaw",
				"err", err1,
				"pidRaw", pidRaw)
			return nil, err1
		}
		if memberId == pid {
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (s *Service) AddParticipant(ctx context.Context, conversationId string, memberId uuid.UUID) error {
	err := s.repository.AddParticipantId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveParticipant(ctx context.Context, conversationId string, memberId uuid.UUID) error {
	err := s.repository.RemoveParticipantId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
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
	c, err := s.repository.FindConversation(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	var isRegistrant bool
	for _, r := range c.RegistrantIds {
		if bytes.Equal(r.Data, memberId[:]) {
			isRegistrant = true
			break
		}
	}
	modIds := make([]uuid.UUID, 0, len(c.ModeratorIds))
	for _, m := range c.ModeratorIds {
		modIds = append(modIds, uuid.UUID(m.Data))
	}

	canEnter := true
	if time.Now().UTC().Before(c.Time.Add(-15 * time.Minute)) {
		canEnter = false
	}
	if time.Now().UTC().Before(c.Time.Add(10*time.Minute)) && !isRegistrant {
		canEnter = false
	}
	for _, b := range c.BanIds {
		if bytes.Equal(b.Data, memberId[:]) {
			canEnter = false
			break
		}
	}
	resp := dto.OnlineConversationDetailResponse{
		Id:           uuid.UUID(c.Id.Data),
		Novel:        c.Novel,
		ShortStory:   c.ShortStory,
		Poem:         c.Poem,
		Play:         c.Play,
		Film:         c.Film,
		WrittenBy:    c.WrittenBy,
		Rule:         c.Rule,
		Capacity:     c.Capacity,
		Time:         c.Time,
		Length:       c.Length.String(),
		CanEnter:     canEnter,
		ModeratorIds: modIds,
		IsRegistrant: isRegistrant,
	}
	return &resp, nil
}

func (s *Service) BanParticipant(ctx context.Context, modId, conversationId, banId uuid.UUID) error {
	mIds, err := s.repository.FindModeratorIds(ctx, conversationId)
	if err != nil {
		return err
	}
	isMod := false
	for _, mId := range mIds {
		if bytes.Equal(mId.Data, modId[:]) {
			isMod = true
			break
		}
	}
	if !isMod {
		return errors.New("you cannot ban")
	}
	err = s.repository.AddBanId(ctx, conversationId, banId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ReportOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	ids, err := s.repository.FindReporterIds(ctx, conversationId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if bytes.Equal(id.Data, memberId[:]) {
			return nil
		}
	}
	if len(ids) > 5 {
		err = s.repository.DeleteOnlineConversation(ctx, conversationId)
		if err != nil {
			return err
		}
	}
	err = s.repository.AddReporterId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RegisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	capacity, err := s.repository.FindCapacity(ctx, conversationId)
	if err != nil {
		return err
	}
	err = s.repository.AddRegistrantId(ctx, conversationId, memberId, capacity)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeregisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	err := s.repository.RemoveRegistrantId(ctx, conversationId, memberId)
	if err != nil {
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
	c, err := s.repository.FindConversation(ctx, conversationId)
	if err != nil {
		return err
	}
	err = s.repository.AddNotificationId(ctx, conversationId, memberId)
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
		p.ScheduledTime = c.Time.Add(-15 * time.Minute).Unix()
		p.Contents = map[int]string{0: string(aboutRaw)}
		p.Type = "online-conversation"
	}
	s.producer.PushMessage("scheduled-notification", nil,
		payload.Marshal(p),
		nil)
	return nil
}

func (s *Service) CancelNotification(ctx context.Context, memberId, conversationId uuid.UUID) error {
	err := s.repository.RemoveNotificationId(ctx, conversationId, memberId)
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
