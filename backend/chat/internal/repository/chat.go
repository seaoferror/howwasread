package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) SaveLike(ctx context.Context, fromId, toId uuid.UUID) error {
	return nil
}
