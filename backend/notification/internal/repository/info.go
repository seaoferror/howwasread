package repository

import (
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) FindNameById(ctx context.Context, id gocql.UUID) (name string, err error) {
	err = r.session.Query("SELECT name FROM profile_by_id WHERE id = ?", id).ScanContext(ctx, &name)
	if err != nil {
		slog.Error("fail to select name from profile_by_id",
			"err", err,
			"id", id.String())
		return "", err
	}
	return name, nil
}

func (r *Repository) FindRoomNameById(ctx context.Context, id gocql.UUID) (name string, err error) {
	err = r.session.Query("SELECT name FROM chat_room_by_id", id).ScanContext(ctx, name)
	if err != nil {
		slog.Error("fail to select name from chat_room_by_id",
			"err", err,
			"id", id.String())
		return "", err
	}
	return name, nil
}
