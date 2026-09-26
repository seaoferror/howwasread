package service

import (
	"backend/common"
	"backend/messagerelay/internal/client"
	"backend/messagerelay/internal/repository"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	RelayMessage(ctx context.Context, id uuid.UUID, toIds [][]byte, roomId, fromId uuid.UUID, contentType string, contents []string)
}

type service struct {
	repository  repository.Repository
	producer    common.Producer
	relayClient client.RelayClient
}

func NewService(r repository.Repository, p common.Producer, relayClient client.RelayClient) Service {
	return &service{
		repository:  r,
		producer:    p,
		relayClient: relayClient,
	}
}
