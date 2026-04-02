package service

import (
	"backend/payload"
	pb "backend/proto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) SendLike(ctx context.Context, fromId, toId uuid.UUID) error {
	messageId, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid v7 for saving payload", "err", err)
		return err
	}
	err = s.repository.SaveLike(ctx, gocql.UUID(messageId), gocql.UUID(fromId), gocql.UUID(toId))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) PublishPersonalMessaging(toId, fromId uuid.UUID, contentType, content string) error {
	err := s.publishMessaging([]uuid.UUID{toId}, []byte{}, fromId, contentType, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) publishMessaging(
	toIds []uuid.UUID,
	roomId []byte,
	fromId uuid.UUID,
	contentType,
	content string) error {
	p := payload.ChatMessaging{
		FromId:      fromId[:],
		RoomId:      roomId,
		ContentType: contentType,
		Content:     content,
	}

	for _, toId := range toIds {
		p.ToIds = append(p.ToIds, toId[:])
	}

	pRaw, err := json.Marshal(p)
	if err != nil {
		slog.Error("fail to Marshal", "err", err)
		return err
	}
	err = s.kafkaProducer.PushMessage("chat.messaging", pRaw)
	if err != nil {
		slog.Error("fail to publish message", "err", err)
		return err
	}
	return nil
}

func (s *Service) NotifyMessaging(
	ctx context.Context,
	toIds [][]byte, roomId, fromId uuid.UUID,
	contentType, content string) error {
	ec := make(chan error)
	var wg sync.WaitGroup
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.repository.RemoveServerIP(ctx, string(toId[:]))
			if err != nil {
				ec <- err
				return
			}
		}()
	}
	wg.Wait()

	client := pb.NewNotificationServiceClient(s.clientConn)
	req := pb.NotifyMessagingRequest{
		ToIds:       toIds,
		FromId:      fromId[:],
		ContentType: contentType,
		Content:     content,
	}
	if roomId != uuid.Nil {
		req.RoomId = roomId[:]
	}
	ctxt, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := client.NotifyMessaging(ctxt, &req)
	if err != nil {
		slog.Error("fail to relay messaging",
			"err", err)
		ec <- err
	}
	close(ec)
	var es []error
	for e := range ec {
		es = append(es, e)
	}

	return errors.Join(es...)
}
