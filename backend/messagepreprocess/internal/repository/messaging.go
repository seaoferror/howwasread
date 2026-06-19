package repository

import (
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) FindParticipantIds(ctx context.Context, roomId gocql.UUID) (participantIds []gocql.UUID, err error) {
	err = r.session.Query(`SELECT participant_ids FROM chat_room_by_id WHERE id = ?`,
		roomId).ScanContext(ctx, &participantIds)
	if err != nil {
		slog.Error("fail to find participant ids by room id", "err", err)
		return nil, err
	}
	return participantIds, nil
}

func (r *Repository) SaveIdsByFileName(ctx context.Context, ids []gocql.UUID, filename gocql.UUID) error {
	err := r.session.Query(`INSERT INTO ids_by_filename (ids, filename) VALUES (?, ?)`,
		ids, filename).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to ids by filename", "err", err)
		return err
	}
	return nil
}

func (r *Repository) CreateChatRoom(ctx context.Context, roomId gocql.UUID, memberId gocql.UUID, roomName string) error {
	err := r.session.Query(`INSERT INTO chat_room_by_id (id, name, room_type, participant_ids) VALUES (?, ?, ?, ?)`,
		roomId, roomName, "group", []gocql.UUID{memberId}).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to create chat room", "err", err,
			"roomId", roomId,
			"memberId", memberId)
		return err
	}
	return nil
}

func (r *Repository) AddParticipantId(ctx context.Context, roomId gocql.UUID, participantId gocql.UUID) error {
	err := r.session.Query(`UPDATE chat_room_by_id SET participant_ids = participant_ids + ? WHERE id = ?`,
		participantId, roomId).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to save participant id", "err", err,
			"roomId", roomId,
			"participantId", participantId)
		return err
	}
	return nil
}

func (r *Repository) RemoveParticipantId(ctx context.Context, roomId gocql.UUID, participantId gocql.UUID) error {
	err := r.session.Query(`UPDATE chat_room_by_id SET participant_ids = participant_ids - ? WHERE id = ?`,
		participantId, roomId).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to delete participant id", "err", err,
			"roomId", roomId,
			"participantId", participantId)
		return err
	}
	return nil
}
