package service

import (
	"backend/onlineconversation/internal/dto"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (s *Service) CreateConversation(ctx context.Context, memberId uuid.UUID, req dto.CreateConversationRequest) (map[string]uuid.UUID, error) {
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

	err = s.repository.InsertConversation(ctx, tx,
		conversationId,
		req)
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
	return map[string]uuid.UUID{"conversationId": conversationId}, nil
}

func (s *Service) UpdateConversation(ctx context.Context, memberId uuid.UUID, req dto.UpdateConversationRequest) error {
	ok, err := s.repository.UpdateConversationIfModerator(ctx, s.repository.Tx(), memberId, req)
	if err != nil {
		return err
	}
	if !ok {
		err = errors.New("can't update conversation")
		slog.Warn("update failed, ui error or api abuse attempt",
			"conversationId", req.Id,
			"memberId", memberId)
		return err
	}
	return nil
}

func (s *Service) DeleteConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	ok, err := s.repository.DeleteOnlineConversationIfModerator(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return err
	}
	if !ok {
		err = errors.New("can't delete conversation")
		slog.Warn("delete failed, ui error or api abuse attempt",
			"conversationId", conversationId,
			"memberId", memberId)
		return err
	}
	return nil
}

func (s *Service) BanParticipant(ctx context.Context, modId, conversationId, banId uuid.UUID) error {
	ok, err := s.repository.AddBanIdIfModerator(ctx, s.repository.Tx(), conversationId, modId, banId)
	if err != nil {
		return err
	}
	if !ok {
		err = errors.New("can't ban participant")
		slog.Warn("ban failed, ui error, or api abuse attempt",
			"conversationId", conversationId,
			"modId", modId,
			"banId", banId)
		return err
	}
	return nil
}

func (s *Service) GetConversationDetail(ctx context.Context, conversationId, memberId uuid.UUID) (*dto.OnlineConversationDetailResponse, error) {
	c, isModerator, isRegistrant, isBanned, isNotificationScheduled, err := s.repository.FindConversationDetail(ctx, s.repository.Tx(), conversationId, memberId)
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
		LengthMinutes:           c.LengthMinutes,
		CanEnter:                canEnter,
		IsModerator:             isModerator,
		IsRegistrant:            isRegistrant,
		IsNotificationScheduled: isNotificationScheduled,
	}
	return &resp, nil
}

func (s *Service) ReportOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	alreadyReported, count, err := s.repository.FindReportStatus(ctx, tx, conversationId, memberId)
	if err != nil {
		return err
	}
	if alreadyReported {
		return nil
	}
	if count > 5 {
		err = s.repository.DeleteOnlineConversation(ctx, tx, conversationId)
		if err != nil {
			return err
		}
	} else {
		err = s.repository.AddReporterId(ctx, tx, conversationId, memberId)
		if err != nil {
			return err
		}
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
	err := s.repository.AddNotificationId(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CancelNotification(ctx context.Context, memberId, conversationId uuid.UUID) error {
	err := s.repository.RemoveNotificationId(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}
