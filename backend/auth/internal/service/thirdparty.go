package service

import (
	"backend/auth/internal/dto"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) provideSessionId(email string) (*dto.SignInWithThirdPartyResponse, string, error) {
	sessionId := uuid.New()
	err := s.repository.SaveEmailBySessionId(gocql.UUID(sessionId), email)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	resp := dto.SignInWithThirdPartyResponse{
		SessionId: sessionId,
	}
	return &resp, "", nil
}

func (s *Service) provideTokens(id gocql.UUID, role string) (*dto.SignInWithThirdPartyResponse, string, error) {
	jti, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to create random uuid for jti")
		return nil, "", ErrInternalServer
	}
	at, rt, err := s.createLoginTokens(id.String(), jti.String(), role)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	err = s.repository.SaveRefreshTokenJTIById(id, jti)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	resp := dto.SignInWithThirdPartyResponse{
		AccessToken: at,
	}
	return &resp, rt, nil
}
