package service

import (
	"backend/common"
	"backend/signalrelay/internal/client"
	"backend/signalrelay/internal/repository"
	"context"
	"encoding/json"
)

type Service interface {
	PropagateSignal(ctx context.Context, toIds [][]byte, fromId []byte, signal json.RawMessage)
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
