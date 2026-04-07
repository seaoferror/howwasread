package repository

import (
	"backend/submessaging/internal/constant"
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) FindParticipantIds(ctx context.Context, roomId gocql.UUID) (participantIds []gocql.UUID, err error) {
	err = r.session.Query(`SELECT participant_ids FROM room_by_id WHERE id = ?`,
		roomId).ScanContext(ctx, &participantIds)
	if err != nil {
		slog.Error("fail to find participant ids by room id", "err", err)
		return nil, err
	}
	return participantIds, nil
}

func (r *Repository) SaveMessaging(ctx context.Context, id, toId, fromId gocql.UUID, roomId []byte, contentType, content string) error {
	err := r.session.Query(`INSERT INTO message_by_to_id 
    (id, to_id, from_id, room_id, content_type, content) VALUES (?, ?, ?, ?, ?, ?) USING TTL ?`,
		id, toId, fromId, roomId, contentType, content, constant.Message_TTL).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to save messaging by toId", "err", err)
		return err
	}
	return nil
}
