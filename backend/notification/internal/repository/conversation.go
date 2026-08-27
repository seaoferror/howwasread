package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) SaveConversationNotification(ctx context.Context, t time.Time, conversationType string, conversationId, memberId uuid.UUID) error {
	err := r.session.Query(
		`INSERT INTO conversation_notification_by_scheduled_time (scheduled_time, conversation_type, conversation_id, member_id) VALUES (?, ?, ?, ?)`,
		t, conversationType, conversationId, memberId).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to save conversation notification",
			"err", err,
			"time", t,
			"conversationType", conversationType,
			"conversationId", conversationId,
			"memberId", memberId)
		return err
	}
	return nil
}
