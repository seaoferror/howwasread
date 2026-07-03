package service

import (
	"backend/common/payload"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) SendNotification(
	ctx context.Context,
	//originTopic string,
	//retryBackoff time.Duration,
	messageId uuid.UUID,
	notificationId uint8,
	value []byte) {
	did, err := s.repository.DidNotification(ctx, string(messageId[:]), string(notificationId))
	if err != nil {
		//err = s.producer.PushRetryMessage(
		//	fmt.Sprintf("%v-%v", originTopic, "retry"),
		//	append(messageId[:], notificationId), value,
		//	50*time.Millisecond,
		//	err.Error(),
		//)
		//if err != nil {
		//	panic(err)
		//}
		return
	}
	if did {
		return
	}
	var p payload.NotificationMessage
	err = json.Unmarshal(value, &p)
	if err != nil {
		slog.Error("fail to unmarshal payload value",
			"err", err,
			"value", value)
		return
	}
	tokenMap := p.TokenMap
	roomName := p.RoomName
	senderName := p.SenderName
	text := p.Text
	imageURL := p.ImageURL
	var wg sync.WaitGroup
	var em sync.Mutex
	var es []error
	var ts []string
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
		return
	}
	err = s.repository.MarkNotification(ctx, string(messageId[:]), string(notificationId))
	if err != nil {
		return
	}
	if br.FailureCount > 0 {
		for i, resp := range br.Responses {
			if messaging.IsUnregistered(resp.Error) || messaging.IsInvalidArgument(resp.Error) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err2 := s.repository.RemoveNotificationInfoByIdAndToken(
						ctx, gocql.UUID(tokenMap[ts[i]]), ts[i])
					if err2 != nil {
						em.Lock()
						es = append(es, err2)
						em.Unlock()
					}
				}()
			}
		}
	}
	wg.Wait()
	err = errors.Join(es...)
	if err != nil {
		return
	}
	return
}
