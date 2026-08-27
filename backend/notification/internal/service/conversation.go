package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (s *Service) ScheduleOnlineConversationNotification(
	ctx context.Context,
	conversationId,
	memberId uuid.UUID,
	t time.Time,
) {
	err := s.repository.SaveConversationNotification(ctx, t, "online", conversationId, memberId)
	if err != nil {
		return
	}
}

func (s *Service) executeConversationNotifications(ctx context.Context) {
	var lastProcessedBucket time.Time
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Info("leadership context canceled, stopping cron polling")
			return
		case now := <-ticker.C:
			currentBucket := now.UTC().Truncate(time.Minute)
			if currentBucket.After(lastProcessedBucket) {
				slog.Info("processing new time bucket", "bucket", currentBucket)

				// 1. Fetch from Cassandra using your repository
				// Example: pendingIDs, err := s.repository.GetPendingConversations(ctx, currentBucket)

				// 2. Fan-out to Kafka using your producer
				// Example:
				// for _, id := range pendingIDs {
				//     s.producer.PublishNotificationTask(id)
				// }

				// 3. (Optional but recommended) Mark the bucket as processed in Cassandra
				// Example: s.repository.MarkBucketComplete(ctx, currentBucket)

				// Update the tracker so this bucket is ignored on the next 5-second tick
				lastProcessedBucket = currentBucket
			}
		}
	}
}
