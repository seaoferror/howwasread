package service

import (
	"backend/common"
	"backend/fcmnotification/internal/client"
	"backend/fcmnotification/internal/repository"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	SendNotification(ctx context.Context, messageId uuid.UUID, notificationId uint8, value []byte)
}

type service struct {
	producer   common.Producer
	repository repository.Repository
	fcmClient  client.FCMClient
}

func NewService(r repository.Repository, p common.Producer, fcmClient client.FCMClient) Service {
	s := service{
		producer:   p,
		repository: r,
		fcmClient:  fcmClient,
	}
	return &s
}
