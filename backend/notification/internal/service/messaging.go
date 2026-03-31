package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) NotifyMessaging(ctx context.Context,
	toIds []uuid.UUID, roomId, fromId uuid.UUID,
	contentType, content string) error {

	var wg sync.WaitGroup
	ec := make(chan error)
	tc := make(chan map[uuid.UUID]string, len(toIds))
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			t, err := s.repository.FindDevicePushToken(ctxt, gocql.UUID(toId))
			if err != nil {
				ec <- err
				return
			}
			tc <- map[uuid.UUID]string{toId: t}
		}()
	}
	wg.Wait()
	close(tc)
	tm := make(map[string]uuid.UUID)
	var tokens []string
	for t := range tc {
		for k, v := range t {
			tm[v] = k
			tokens = append(tokens, v)
		}
	}

	func() {
		client, err := s.app.Messaging(ctx)
		if err != nil {
			slog.Error("fail to create fcm client",
				"err", err)
			ec <- err
			return
		}
		data := make(map[string]string)
		if roomId != uuid.Nil {
			data["roomId"] = roomId.String()
		}
		data["fromId"] = fromId.String()
		data["contentType"] = contentType
		data["content"] = content

		message := &messaging.MulticastMessage{
			Data:   data,
			Tokens: tokens,
		}
		br, err := client.SendEachForMulticast(ctx, message)
		if err != nil {
			ec <- err
			return
		}
		if br.FailureCount > 0 {
			var failedIds []string
			for i, resp := range br.Responses {
				if !resp.Success {
					// The order of responses corresponds to the order of the registration tokens.
					failedIds = append(failedIds, tm[tokens[i]].String())
				}
			}
			slog.Error("failed ids to get push notification", "failedIds", failedIds)
			//TODO: publish this to kafka and retry?
		}
	}()

	close(ec)
	var errs []error
	for err := range ec {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

//I can trigger sqlite message saving by fcm, so I don't need to save any message in my database
//for _, toId := range toIds {
//	wg.Add(1)
//	go func(capturedToId uuid.UUID) {
//		defer wg.Done()
//		ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
//		defer cancel()
//		err := s.repository.SaveMessaging(
//			ctxt, gocql.UUID(capturedToId), gocql.UUID(roomId), gocql.UUID(fromId),
//			contentType, content)
//		if err != nil {
//			ec <- err
//		}
//	}(toId)
//}
//wg.Wait()
