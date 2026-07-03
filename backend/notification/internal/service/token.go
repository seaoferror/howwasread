package service

import (
	"bytes"
	"context"
	"errors"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) RegisterNotification(ctx context.Context, id uuid.UUID, os, token string) error {
	old, err := s.repository.FindMemberIdByToken(ctx, token)
	if errors.Is(err, gocql.ErrNotFound) {
		err = nil
		err = s.repository.UpdateMemberIdByToken(ctx, token, id)
		if err != nil {
			return err
		}
	}
	if err != nil {
		slog.Error("fail to find member Id by Token", "err", err,
			"token", token)
		return err
	}
	if !bytes.Equal(old[:], id[:]) {
		err = s.repository.DeleteNotificationInfoByIdAndToken(ctx, old, token)
		if err != nil {
			return err
		}
		err = s.repository.UpdateMemberIdByToken(ctx, token, id)
		if err != nil {
			return err
		}
	}
	err = s.repository.SaveNotificationInfoById(ctx, gocql.UUID(id), os, token)
	if err != nil {
		return err
	}
	return nil
}
