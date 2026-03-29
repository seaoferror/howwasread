package repository

import (
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) SaveLike(ctx context.Context, messageId, fromId, toId gocql.UUID) error {
	err := r.session.Query(`INSERT INTO message_by_to_id (
                     id, to_id_type, to_id, from_id, content_type, content, read) VALUES (?, ?, ?, ?, ?)`,
		messageId, "personal", toId, fromId, "text", "👍", false).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to save like", "err", err)
		return err
	}
	return nil
}
