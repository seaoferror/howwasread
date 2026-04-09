package repository

import (
	"backend/notification/internal/data"
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) SetPushTokenById(ctx context.Context, id, deviceId gocql.UUID, os, token string) error {
	err := r.session.Query(`INSERT INTO device_push_token_by_id (id, deviceId, os, token) VALUES (?, ?, ?, ?)`,
		id, deviceId, os, token).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to set device push token",
			"err", err)
		return err
	}
	return nil
}

func (r *Repository) FindPushTokensById(ctx context.Context, id gocql.UUID) (result []data.FindPushTokensById, err error) {
	iter := r.session.Query(`SELECT device_id, os, token FROM device_push_token_by_id WHERE id = ?`,
		id).IterContext(ctx)
	var deviceId gocql.UUID
	var os, token string
	for iter.Scan(&deviceId, &os, &token) {
		result = append(result, data.FindPushTokensById{
			Id:              id,
			DeviceId:        deviceId,
			OS:              os,
			DevicePushToken: token,
		})
	}
	err = iter.Close()
	if err != nil {
		slog.Error("fail to close iterator", "err", err)
		return nil, err
	}
	return result, nil
}

func (r *Repository) RemovePushTokenByIdAndDeviceId(ctx context.Context, id, deviceId gocql.UUID) error {
	err := r.session.Query(`DELETE FROM device_push_token_by_id WHERE id = ? AND device_id = ?`,
		id, deviceId).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to remove invalid device token", "err", err)
		return err
	}
	return nil
}
