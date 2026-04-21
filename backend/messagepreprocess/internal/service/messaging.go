package service

import (
	"backend/payload"
	"context"
	"encoding/json"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) ManageMessage(ctx context.Context, id uuid.UUID, fromId uuid.UUID, toIdType string, toId uuid.UUID, contentType string, content string) error {
	var toIds [][]byte
	roomId := toId
	if toIdType == "personal" {
		toIds = append(toIds, toId[:])
		toIds = append(toIds, fromId[:])
	}
	if toIdType == "room" {
		participantIds, err := s.repository.FindParticipantIds(ctx, gocql.UUID(toId))
		if err != nil {
			return err
		}
		for _, pid := range participantIds {
			toIds = append(toIds, pid[:])
		}
	}
	if contentType != "text" {
		filename, err := uuid.Parse(content)
		if err != nil {
			slog.Error("fail to parse uuid from content", "err", err)
			return err
		}
		var castedIds []gocql.UUID
		for _, tid := range toIds {
			castedIds = append(castedIds, gocql.UUID(tid))
		}
		err = s.repository.SaveIdsByFileName(ctx, castedIds, gocql.UUID(filename))
		if err != nil {
			return err
		}
	}

	p, _ := json.Marshal(payload.PreparedMessage{
		Id:          id[:],
		ToIds:       toIds,
		RoomId:      roomId[:],
		FromId:      fromId[:],
		ContentType: contentType,
		Content:     content,
	})

	err := s.producer.PushMessage("manage_message.prepared", nil, p)
	if err != nil {
		return err
	}

	return nil
}
