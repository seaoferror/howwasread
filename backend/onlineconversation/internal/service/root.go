package service

import (
	"backend/common"
	"backend/onlineconversation/internal/dto"
	"backend/onlineconversation/internal/repository"
	"context"
	"os"

	"github.com/google/uuid"
)

type Service interface {
	GenerateTurn() *dto.GetTurnResponse
	GetParticipantsWithoutMe(ctx context.Context, conversationId string, memberId uuid.UUID) ([]uuid.UUID, error)
	AddParticipant(ctx context.Context, conversationId string, memberId uuid.UUID) error
	RemoveParticipant(ctx context.Context, conversationId string, memberId uuid.UUID) error
	SetServerIP(ctx context.Context, memberId uuid.UUID, ip string) error
	RemoveServerIP(ctx context.Context, memberId uuid.UUID) error
	PublishConversationSignal(fromId uuid.UUID, toIds [][]byte, signal []byte) error
	CreateConversation(ctx context.Context, memberId uuid.UUID, req dto.CreateConversationRequest) (map[string]uuid.UUID, error)
	UpdateConversation(ctx context.Context, memberId uuid.UUID, req dto.UpdateConversationRequest) error
	DeleteConversation(ctx context.Context, memberId, conversationId uuid.UUID) error
	BanParticipant(ctx context.Context, modId, conversationId, banId uuid.UUID) error
	GetConversationDetail(ctx context.Context, conversationId, memberId uuid.UUID) (*dto.OnlineConversationDetailResponse, error)
	RegisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error
	DeregisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error
	ScheduleNotification(ctx context.Context, memberId, conversationId uuid.UUID) error
	CancelNotification(ctx context.Context, memberId, conversationId uuid.UUID) error
}

type service struct {
	repository repository.Repository
	producer   common.Producer
	turnSecret string
	turnRealm  string
}

func NewService(r repository.Repository, p common.Producer) Service {
	s := &service{
		repository: r,
		producer:   p,
		turnSecret: os.Getenv("TURN_SECRET"),
		turnRealm:  os.Getenv("TURN_REALM"),
	}

	return s
}
