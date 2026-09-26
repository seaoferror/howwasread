package service

import (
	"backend/common"
	"backend/notification/internal/repository"
	"context"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	PreprocessMessageNotification(ctx context.Context, notificationId uint8, messageId uuid.UUID, toIds [][]byte, roomId, fromId uuid.UUID, contentType string, content []string)
	PreprocessScheduledNotification(ctx context.Context, partitionId uuid.UUID, notifications map[uuid.UUID]map[int]string, contents map[int]string)
	RegisterNotification(ctx context.Context, id uuid.UUID, os, token string) error
}

type service struct {
	repository repository.Repository
	producer   common.Producer
	cdnClient  common.CDNClient
}

func NewService(r repository.Repository, p common.Producer, cdnClient common.CDNClient) Service {
	s := service{
		repository: r,
		producer:   p,
		cdnClient:  cdnClient,
	}

	return &s
}
