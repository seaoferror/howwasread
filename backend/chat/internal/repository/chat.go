package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *Repository) SaveLike(ctx context.Context, messageId, fromId, toId uuid.UUID) error {

	_, err := r.db.ExecContext(ctx, `INSERT INTO message (
                     id, to_id, from_id, content_type, content) VALUES (?, ?, ?, ?, ?)`,
		messageId[:], toId[:], fromId[:], "text", "👍")
	if err != nil {
		slog.Error("fail to save like", "err", err)
		return err
	}
	return nil
}
