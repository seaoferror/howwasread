package service

import (
	"backend/payload"
	"context"
	"encoding/json"
	"log/slog"

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
