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
		slog.Warn("update online conversation failed, ui error or api abuse attempt",
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
		slog.Warn("delete online conversation failed, ui error or api abuse attempt",
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
	detail, err := s.repository.FindConversationDetail(ctx, s.repository.Tx(), conversationId, memberId)
	if err != nil {
		return nil, err
	}
	canEnter := true
	if time.Now().UTC().Before(detail.Time.Add(-15 * time.Minute)) {
		canEnter = false
	}
	if time.Now().UTC().Before(detail.Time.Add(10*time.Minute)) && !detail.IsRegistrant {
		canEnter = false
	}
	if detail.IsBanned {
		canEnter = false
	}
	resp := dto.OnlineConversationDetailResponse{
		Novel:                   detail.Novel,
		ShortStory:              detail.ShortStory,
		Poem:                    detail.Poem,
		Play:                    detail.Play,
		Film:                    detail.Film,
		WrittenBy:               detail.WrittenBy,
		Rule:                    detail.Rule,
		Capacity:                detail.Capacity,
		Time:                    detail.Time,
		LengthMinutes:           detail.LengthMinutes,
		CanEnter:                canEnter,
		IsModerator:             detail.IsModerator,
		IsRegistrant:            detail.IsRegistrant,
		IsNotificationScheduled: detail.IsNotificationScheduled,
	}
	return &resp, nil
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
