package service

import (
	"backend/payload"
	pb "backend/proto"
	"context"
	"encoding/json"
	"log/slog"
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

func (s *Service) PublishMessaging(id, fromId uuid.UUID, toIdType string, toId uuid.UUID, contentType, content string) error {
	p, _ := json.Marshal(payload.ChatMessaging{
		Id:          id[:],
		FromId:      fromId[:],
		ToIdType:    toIdType,
		ToId:        toId[:],
		ContentType: contentType,
		Content:     content,
	})
	err := s.kafkaProducer.PushMessage("chat.messaging", p)
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
		slog.Error("fail to notify messaging",
			"err", err)
		return err
	}
	return nil
}
