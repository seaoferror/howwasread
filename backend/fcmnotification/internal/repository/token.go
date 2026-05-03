package repository

import (
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) RemovePushTokenByIdAndDeviceId(ctx context.Context, id, deviceId gocql.UUID) error {
	err := r.session.Query(`DELETE FROM device_push_token_by_id WHERE id = ? AND device_id = ?`,
		id, deviceId).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to remove invalid device token", "err", err)
		return err
	}
	return nil
}
