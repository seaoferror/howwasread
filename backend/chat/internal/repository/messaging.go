package repository

import (
	"backend/chat/internal/data"
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) FindMessagesByToIdAndId(ctx context.Context, id, cursor gocql.UUID) (result []data.FindMessagesByToIdAndId, err error) {
	iter := r.session.Query(`SELECT id, room_id, from_id, content_type, content FROM message_by_to_id 
                                                   WHERE to_id = ? AND id > ?`,
		id, cursor).IterContext(ctx)
	var messageId, fromId gocql.UUID
	var roomId []byte
	var content, contentType string
	for iter.Scan(&messageId, &roomId, &fromId, &content, &contentType) {
		result = append(result, data.FindMessagesByToIdAndId{
			Id:          messageId,
			RoomId:      roomId,
			FromId:      fromId,
			Content:     content,
			ContentType: contentType,
		})
	}
	err = iter.Close()
	if err != nil {
		slog.Error("fail to close iterator", "err", err)
		return nil, err
	}
	return result, err
}
