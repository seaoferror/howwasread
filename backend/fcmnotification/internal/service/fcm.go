package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) SendNotification(ctx context.Context, messageId uuid.UUID, notificationId uint8, tokenMap map[string]uuid.UUIDs, roomName, senderName, text, imageURL string) error {
	did, err := s.repository.DidNotification(ctx, string(messageId[:]), string(notificationId))
	if err != nil {
		return err
	}
	if did {
		return nil
	}
	var wg sync.WaitGroup
	var em sync.Mutex
	var es []error
	var ts []string
	var failedIds []uuid.UUIDs
	for t := range tokenMap {
		ts = append(ts, t)
	}
	title := senderName
	if roomName != "" {
		title = roomName
		text = fmt.Sprintf("%v: %v", senderName, text)
	}
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     text,
			ImageURL: imageURL,
		},
		Tokens: ts,
	}
	ctxf, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	br, err := s.fcmClient.SendEachForMulticast(ctxf, message)
	if err != nil {
		return err
	}
	err = s.repository.MarkNotification(ctx, string(messageId[:]), string(notificationId))
	if err != nil {
		return err
	}
	if br.FailureCount > 0 {
		for i, resp := range br.Responses {
			if !resp.Success {
				slog.Info("fail to send fcm notification",
					"failedId", tokenMap[ts[i]])
				failedIds = append(failedIds, tokenMap[ts[i]])
			}
		}
	}
	for _, failedId := range failedIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err2 := s.repository.RemovePushTokenByIdAndDeviceId(
				ctx, gocql.UUID(failedId[0]), gocql.UUID(failedId[1]))
			if err2 != nil {
				em.Lock()
				es = append(es, err2)
				em.Unlock()
			}
		}()
	}
	wg.Wait()
	err = errors.Join(es...)
	if err != nil {
		return err
	}
	return nil
}
