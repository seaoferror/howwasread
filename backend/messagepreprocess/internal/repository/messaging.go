package repository

import (
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

func (r *Repository) SaveIdsByFileName(ctx context.Context, ids []gocql.UUID, filename gocql.UUID) error {
	err := r.session.Query(`INSERT INTO ids_by_filename (ids, filename) VALUES (?, ?)`,
		ids, filename).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to ids by filename", "err", err)
		return err
	}
	return nil
}
