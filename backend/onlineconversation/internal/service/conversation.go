package service

import (
	"backend/common/payload"
	"backend/onlineconversation/internal/dto"
	"context"
	"errors"
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
		return nil, err
	}
	defer tx.Rollback()

	err = s.repository.InsertConversation(ctx, tx, conversationId, novel, shortStory, poem, play, film, writtenBy, rule, capacity, t, length)
	if err != nil {
		return nil, err
	}
	err = s.repository.InsertModerator(ctx, tx, conversationId, memberId)
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

func (s *Service) DeleteConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	isModerator, err := s.repository.IsModerator(ctx, tx, conversationId, memberId)
	if err != nil {
		return err
	}
	if !isModerator {
		return errors.New("only moderator can delete conversation")
	}
	err = s.repository.DeleteOnlineConversation(ctx, tx, conversationId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("fail to commit transaction for delete conversation", "err", err)
		return err
	}
	return nil
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

func (s *Service) GetConversationDetail(ctx context.Context, conversationId, memberId uuid.UUID) (*dto.OnlineConversationDetailResponse, error) {
	c, isRegistrant, isBanned, isNotificationScheduled, err := s.repository.FindConversationDetail(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return nil, err
	}
	moderatorIds, err := s.repository.FindModeratorIds(ctx, s.repository.Tx(), conversationId)
	if err != nil {
		return nil, err
	}
	canEnter := true
	if time.Now().UTC().Before(c.Time.Add(-15 * time.Minute)) {
		canEnter = false
	}
	if time.Now().UTC().Before(c.Time.Add(10*time.Minute)) && !isRegistrant {
		canEnter = false
	}
	if isBanned {
		canEnter = false
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
		ModeratorIds:            moderatorIds,
		IsRegistrant:            isRegistrant,
		IsNotificationScheduled: isNotificationScheduled,
	}
	return &resp, nil
}

func (s *Service) BanParticipant(ctx context.Context, modId, conversationId, banId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	isMod, err := s.repository.IsModerator(ctx, tx, conversationId, modId)
	if err != nil {
		return err
	}
	if !isMod {
		return errors.New("you cannot ban")
	}
	err = s.repository.AddBanId(ctx, tx, conversationId, banId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("fail to commit", "err", err)
		return err
	}
	return nil
}

func (s *Service) ReportOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ids, err := s.repository.FindReporterIds(ctx, tx, conversationId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if id == memberId {
			return nil
		}
	}
	if len(ids) > 5 {
		err = s.repository.DeleteOnlineConversation(ctx, tx, conversationId)
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err != nil {
			slog.Error("fail to commit", "err", err)
			return err
		}
		return nil
	}
	err = s.repository.AddReporterId(ctx, tx, conversationId, memberId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("fail to commit", "err", err)
		return err
	}
	return nil
}

func (s *Service) RegisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	registered, err := s.repository.TryIncrementRegistrants(ctx, tx, conversationId)
	if err != nil {
		return err
	}
	if !registered {
		return errors.New("already fully registered")
	}

	err = s.repository.InsertRegistrant(ctx, tx, conversationId, memberId)
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

func (s *Service) ScheduleNotification(ctx context.Context, memberId, conversationId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	c, err := s.repository.FindConversation(ctx, tx, conversationId)
	if err != nil {
		return err
	}
	hasNotification, err := s.repository.HasNotification(ctx, tx, conversationId)
	if err != nil {
		return err
	}
	err = s.repository.AddNotificationId(ctx, tx, conversationId, memberId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("fail to commit transaction", "err", err)
		return err
	}
	p := payload.NotificationScheduling{
		PartitionId: conversationId,
		KeyId:       memberId,
	}
	if !hasNotification {
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
