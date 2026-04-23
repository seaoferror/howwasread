package repository

import (
	"backend/chat/internal/data"
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (r *Repository) FindMessagesByToIdAndId(ctx context.Context, id, cursor gocql.UUID) (result []data.FindMessagesByToIdAndId, err error) {
	iter := r.session.Query(`SELECT id, room_id, from_id, content_type, contents FROM message_by_to_id 
                                                   WHERE to_id = ? AND id > ?`,
		id, cursor).IterContext(ctx)
	var messageId, fromId gocql.UUID
	var roomId []byte
	var contentType string
	var contents []string
	for iter.Scan(&messageId, &roomId, &fromId, &contentType, &contents) {
		result = append(result, data.FindMessagesByToIdAndId{
			Id:          messageId,
			RoomId:      roomId,
			FromId:      fromId,
			ContentType: contentType,
			Contents:    contents,
		})
	}
	err = iter.Close()
	if err != nil {
		slog.Error("fail to close iterator", "err", err)
		return nil, err
	}
	return result, err
}

func (r *Repository) SetFilepath(ctx context.Context, id string, filenames []string) error {
	result := r.client.Do(ctx, r.client.B().Sadd().Key(id).Member(filenames...).Build())
	if result.Error() != nil {
		slog.Error("fail to save member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) HasFilepath(ctx context.Context, id string, filenames []string) (bool, error) {
	for _, f := range filenames {
		result := r.client.Do(ctx, r.client.B().Sismember().Key(id).Member(f).Build())
		if result.Error() != nil {
			slog.Error("fail to check file path", "err", result.Error())
			return false, result.Error()
		}
		exists, err := result.AsBool()
		if err != nil {
			slog.Error("fail to get result of checking file path", "err", err)
			return false, err
		}
		if !exists {
			slog.Info("this filename not exist")
			return false, nil
		}
	}
	return true, nil
}

func (r *Repository) RemoveFilepath(ctx context.Context, id string, filenames []string) error {
	result := r.client.Do(ctx, r.client.B().Srem().Key("presigned"+id).Member(filenames...).Build())
	if result.Error() != nil {
		slog.Error("fail to check file path", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) FindIdsByFilename(ctx context.Context, filename uuid.UUID) (ids []gocql.UUID, err error) {
	err = r.session.Query("SELECT ids FROM ids_by_filename WHERE filename = ?", filename).ScanContext(ctx, &ids)
	if err != nil {
		slog.Error("fail to find ids by filename",
			"err", err)
		return nil, err
	}
	return ids, nil
}
