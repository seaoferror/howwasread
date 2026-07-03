package repository

import (
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) RemoveNotificationInfoByIdAndToken(ctx context.Context, id gocql.UUID, token string) error {
	err := r.session.Query(`DELETE FROM notification_info_by_id WHERE id = ? AND device_push_token = ?`,
		id, token).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to remove invalid device token", "err", err)
		return err
	}
	return nil
}
