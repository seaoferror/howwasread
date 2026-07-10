package service

import (
	"backend/auth/internal/dto"
	"context"
	"errors"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

func (s *Service) SignInWithGoogle(ctx context.Context, token string) (
	*dto.SignInWithThirdPartyResponse,
	string,
	error,
) {
	payload, err := idtoken.Validate(ctx, token, s.googleSignInWebClientId)
	if err != nil {
		return nil, "", ErrSignInWithGoogle
	}
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, "", ErrSignInWithGoogle
	}
	emailVerified, phoneNumberVerified, id, _, role1, err1 := s.repository.FindLoginInfoByEmail(email)
	if errors.Is(err1, gocql.ErrNotFound) {
		err1 = nil
		idv7, err2 := uuid.NewV7()
		if err2 != nil {
			slog.Error("fail to create uuid v7 for google sign in user")
			return nil, "", ErrInternalServer
		}
		err = s.repository.SaveThirdPartySignInInfo(ctx, gocql.UUID(idv7), email, false, true)
		if err != nil {
			return nil, "", ErrInternalServer
		}
		return s.provideSessionId(email)
	}
	if err1 != nil {
		return nil, "", ErrInternalServer
	}
	err = s.repository.SaveThirdPartySignInInfo(ctx, id, email, phoneNumberVerified, emailVerified)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	if !phoneNumberVerified {
		return s.provideSessionId(email)
	}
	return s.provideTokens(id, role1)
}
