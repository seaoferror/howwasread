package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *Service) SendLike(ctx context.Context, fromId, toId uuid.UUID) error {
	messageId, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid v7 for saving message", "err", err)
		return err
	}
	err = s.repository.SaveLike(ctx, messageId, fromId, toId)
	if err != nil {
		return err
	}
	return nil
}
