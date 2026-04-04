package repository

import (
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (r *Repository) SaveMessaging(ctx context.Context, id, toId, roomId, fromId gocql.UUID, contentType, content string) error {
	var rId any
	rId = roomId
	if uuid.UUID(roomId) == uuid.Nil {
		rId = nil
	}
	err := r.session.Query(`INSERT INTO message_by_to_id 
    (id, to_id, room_id, from_id, content_type, content) VALUES (?, ?, ?, ?, ?, ?)`,
		id, toId, rId, fromId, contentType, content).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to save messaging", "err", err)
		return err
	}
	return nil
}

func (r *Repository) SetDevicePushToken(ctx context.Context, id gocql.UUID, os, token string) error {
	err := r.session.Query(`INSERT INTO device_push_token_by_id (id, os, token) VALUES (?, ?, ?)`,
		id, os, token).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to set device push token",
			"err", err)
		return err
	}
	return nil
}

func (r *Repository) FindDevicePushToken(ctx context.Context, id gocql.UUID) (os string, token string, err error) {
	err = r.session.Query(`SELECT os, token FROM device_push_token_by_id WHERE id = ?`,
		id).ScanContext(ctx, &os, &token)
	if err != nil {
		slog.Error("fail to find device push token",
			"err", err)
		return "", "", err
	}
	return os, token, nil
}

func (r *Repository) RemoveInvalidDeviceToken(ctx context.Context, id gocql.UUID) error {
	err := r.session.Query(`UPDATE device_push_token_by_id SET token = null WHERE id = ?`,
		id).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to remove invalid device token", "err", err)
		return err
	}
	return nil
}
