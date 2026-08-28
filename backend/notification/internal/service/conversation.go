package service

import (
	"backend/common/payload"
	"context"
	"fmt"
	"log/slog"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) ScheduleOnlineConversationNotification(
	ctx context.Context,
	t time.Time,
	conversationId,
	memberId uuid.UUID,
	writtenBy string,
) {
	err := s.repository.SaveConversationNotification(ctx, t.Add(-15*time.Minute), "online", writtenBy, conversationId, memberId)
	if err != nil {
		return
	}
}

func (s *Service) executeConversationNotifications(ctx context.Context) {
	var lastProcessedTimeBucket time.Time
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Info("leadership context canceled, stopping cron polling")
			return
		case now := <-ticker.C:
			currentTimeBucket := now.UTC().Truncate(time.Minute)
			if currentTimeBucket.After(lastProcessedTimeBucket) {
				slog.Info("processing new time bucket", "bucket", currentTimeBucket)
				scanner := s.repository.GetConversationNotificationScanner(ctx, currentTimeBucket, "online")
				for scanner.Next() {
					var conversationId gocql.UUID
					var memberIds []gocql.UUID
					var writtenBy string
					err := scanner.Scan(&conversationId, &memberIds, &writtenBy)
					if err != nil {
						slog.Error("fail to scan row of conversation notification", "err", err)
						return
					}
					toIds := make([]uuid.UUID, len(memberIds))
					for _, memberId := range memberIds {
						id := uuid.UUID(memberId)
						toIds = append(toIds, id)
					}
					apntm, fcmtm, err := s.getEachTokenMap(ctx, toIds)
					if err != nil {
						return
					}
					p := payload.NotificationMessage{
						Title1: "Conversation start soon",
						Text:   fmt.Sprintf("Written By: %v", writtenBy),
					}
					p.TokenMap = fcmtm
					s.producer.PushMessage("fcm-notification",
						conversationId[:],
						payload.Marshal(p),
						nil,
					)
					p.TokenMap = apntm
					s.producer.PushMessage("apn-notification",
						conversationId[:],
						payload.Marshal(p),
						nil,
					)
				}
				lastProcessedTimeBucket = currentTimeBucket
			}
		}
	}
}
