package service

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) SetDevicePushToken(ctx context.Context, id, deviceId uuid.UUID, os, token string) error {
	err := s.repository.SetPushTokenById(ctx, gocql.UUID(id), gocql.UUID(deviceId), os, token)
	if err != nil {
		return err
	}
	return nil
}
