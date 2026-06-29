package repository

import (
	"context"
	"errors"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) FindChatRoomInfoById(ctx context.Context, id gocql.UUID) (name string, roomType string, err error) {
	err = r.session.Query("SELECT name, room_type FROM chat_room_by_id WHERE id = ?", id).ScanContext(ctx, &name, &roomType)
	if err != nil {
		slog.Error("fail to find chat room info by id",
			"err", err,
			"id", id.String())
		return "", "", nil
	}
	return name, roomType, nil
}

func (r *Repository) FindProfileById(ctx context.Context, id gocql.UUID) (name string, err error) {
	err = r.session.Query(`SELECT name FROM profile_by_id WHERE id = ?`,
		id).ScanContext(ctx, &name)
	if err != nil {
		slog.Info("fail to find profile by id",
			"err", err,
			"id", id.String())
		return "", err
	}
	return name, nil
}

func (r *Repository) SaveNameById(ctx context.Context, id gocql.UUID, name string) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query("UPDATE profile_by_id SET name = ? WHERE id = ?", name, id).
		Query("INSERT INTO chat_room_by_id (id, name, room_type) VALUES (?, ?, ?)", id, name, "personal").
		ExecContext(ctx)
	if err != nil {
		slog.Error("fail to save profile name by id",
			"err", err,
			"id", id.String(),
			"name", name)
		return err
	}
	return nil
}

func (r *Repository) WasBlocked(ctx context.Context, blockerId gocql.UUID, blockedId gocql.UUID) (bool, error) {
	var a gocql.UUID
	err := r.session.Query(
		`SELECT blocked_id FROM block WHERE blocker_id = ? AND blocked_id = ?`, blockerId, blockedId,
	).ScanContext(ctx, &a)
	if errors.Is(gocql.ErrNotFound, err) {
		return false, nil
	}
	if err != nil {
		slog.Error("fail to check block", "err", err,
			"toId", blockerId,
			"fromId", blockedId)
		return false, err
	}
	return true, nil
}
