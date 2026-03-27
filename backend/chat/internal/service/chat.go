package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) SendLike(ctx context.Context, fromId, toId uuid.UUID) error {
	err := s.repository.SaveLike(ctx, fromId, toId)
}
