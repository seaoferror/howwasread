package service

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) SetDevicePushToken(ctx context.Context, id uuid.UUID, os, token string) error {
	err := s.repository.SetDevicePushToken(ctx, gocql.UUID(id), os, token)
	if err != nil {
		return err
	}
	return nil
}
