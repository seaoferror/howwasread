package service

import (
	"backend/common/payload"
	"context"
	"errors"
	"log/slog"
	"sync"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) ManageMessage(
	ctx context.Context,
	id, fromId uuid.UUID,
	toIdType string,
	toId uuid.UUID,
	contentType string,
	contents []string,
) {
	var toIds [][]byte
	roomId := toId
	if toIdType == "personal" {
		toIds = append(toIds, toId[:])
		toIds = append(toIds, fromId[:])
	}
	if contentType == "quit" {
		err := s.repository.RemoveParticipantId(ctx, gocql.UUID(roomId), gocql.UUID(fromId))
		if err != nil {
			return
		}
	}
	if toIdType == "group" && contentType != "create" {
		participantIds, err := s.repository.FindParticipantIds(ctx, gocql.UUID(toId))
		if errors.Is(err, gocql.ErrNotFound) {
			err = nil
		}
		if err != nil {
			return
		}
		for _, pid := range participantIds {
			toIds = append(toIds, pid[:])
		}
	}
	if contentType == "create" {
		err := s.repository.CreateChatRoom(ctx, gocql.UUID(roomId), gocql.UUID(fromId), contents[0])
		contents = nil
		if err != nil {
			return
		}
		toIds = append(toIds, fromId[:])
	}
	if contentType == "participate" {
		err := s.repository.AddParticipantId(ctx, gocql.UUID(roomId), gocql.UUID(fromId))
		if err != nil {
			return
		}
		toIds = append(toIds, fromId[:])
	}
	if contentType == "image" || contentType == "audio" || contentType == "video" {
		var wg sync.WaitGroup
		var em sync.Mutex
		var es []error
		for _, content := range contents {
			wg.Add(1)
			go func() {
				defer wg.Done()
				filename, err := uuid.Parse(content)
				if err != nil {
					slog.Error("fail to parse uuid from content", "err", err)
					em.Lock()
					es = append(es, err)
					em.Unlock()
				}
				var castedIds []gocql.UUID
				for _, tid := range toIds {
					castedIds = append(castedIds, gocql.UUID(tid))
				}
				err = s.repository.SaveIdsByFileName(ctx, castedIds, gocql.UUID(filename))
				if err != nil {
					em.Lock()
					es = append(es, err)
					em.Unlock()
				}
			}()
		}
		wg.Wait()
		err := errors.Join(es...)
		if err != nil {
			return
		}
	}

	p := payload.Marshal(payload.PreparedMessage{
		Id:          id[:],
		ToIds:       toIds,
		RoomId:      roomId[:],
		FromId:      fromId[:],
		ContentType: contentType,
		Contents:    contents,
	})

	err := s.producer.PushMessage("manage_message.prepared", nil, p)
	if err != nil {
		return
	}
	return
}
