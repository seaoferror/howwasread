package repository

import (
	"context"
	"log/slog"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (r *Repository) SaveConversationNotification(
	ctx context.Context,
	t time.Time,
	conversationType, writtenBy string,
	conversationId, memberId uuid.UUID,
) error {
	err := r.session.Query(`
		UPDATE conversation_notification_by_scheduled_time 
		SET member_ids = member_ids + ?,
		    written_by = ?
		WHERE scheduled_time = ? 
		  AND conversation_type = ? 
		  AND conversation_id = ?`,
		[]uuid.UUID{memberId}, writtenBy, t, conversationType, conversationId).
		ExecContext(ctx)
	if err != nil {
		slog.Error("fail to append member to conversation notification",
			"err", err,
			"time", t,
			"conversationType", conversationType,
			"writtenBy", writtenBy,
			"conversationId", conversationId,
			"memberId", memberId,
		)
		return err
	}
	return nil
}

func (r *Repository) GetConversationNotificationScanner(ctx context.Context, bucket time.Time, conversationType string) gocql.Scanner {
	return r.session.Query(`
		SELECT conversation_id, member_ids, written_by
		FROM conversation_notification_by_scheduled_time 
		WHERE scheduled_time = ? AND conversation_type = ?`,
		bucket, conversationType).
		PageSize(500).
		IterContext(ctx).
		Scanner()
}
