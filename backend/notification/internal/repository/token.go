package repository

import (
	"backend/notification/internal/projection"
	"context"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (r *Repository) SaveNotificationInfoById(ctx context.Context, id gocql.UUID, os, token string) error {
	err := r.session.Query(`INSERT INTO notification_info_by_id (id, os, device_push_token) VALUES (?, ?, ?)`,
		id, os, token).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to set device push token",
			"err", err)
		return err
	}
	return nil
}

func (r *Repository) FindPushTokensById(ctx context.Context, id gocql.UUID) (result []projection.FindPushTokensById, err error) {
	iter := r.session.Query(`SELECT os, device_push_token FROM notification_info_by_id WHERE id = ?`,
		id).IterContext(ctx)
	var os, token string
	for iter.Scan(&os, &token) {
		result = append(result, projection.FindPushTokensById{
			Id:              id,
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

func (r *Repository) FindMemberIdByToken(ctx context.Context, token string) (id uuid.UUID, err error) {
	err = r.session.Query(`SELECT id FROM id_by_device_push_token WHERE device_push_token = ?`, token).
		ScanContext(ctx, &id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *Repository) UpdateMemberIdByToken(ctx context.Context, token string, id uuid.UUID) error {
	err := r.session.Query(`INSERT INTO id_by_device_push_token (id, device_push_token) VALUES (?, ?)`, id, token).
		ExecContext(ctx)
	if err != nil {
		slog.Error("fail to update member id by token", "err", err,
			"token", token,
			"id", id)
		return err
	}
	return nil
}

func (r *Repository) DeleteNotificationInfoByIdAndToken(ctx context.Context, id uuid.UUID, token string) error {
	err := r.session.Query(`DELETE FROM notification_info_by_id WHERE id = ? AND device_push_token = ?`, id, token).
		ExecContext(ctx)
	if err != nil {
		slog.Error("fail to delete notification info by id and token", "err", err,
			"id", id,
			"token", token)
		return err
	}
	return nil
}
