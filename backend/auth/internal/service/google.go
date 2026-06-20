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
	userId := payload.Subject
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, "", ErrSignInWithGoogle
	}
	id, emailFromDB, role, err := s.repository.FindThirdPartySignInInfo(ctx, userId, "google")
	if errors.Is(err, gocql.ErrNotFound) {
		err = nil
		emailVerified, phoneNumberVerified, id1, _, role1, err1 := s.repository.FindLoginInfoByEmail(email)
		if errors.Is(err1, gocql.ErrNotFound) {
			err1 = nil
			idv7, err2 := uuid.NewV7()
			if err2 != nil {
				slog.Error("fail to create uuid v7 for apple sign in user")
				return nil, "", ErrInternalServer
			}
			err = s.repository.SaveThirdPartySignInInfo(ctx, gocql.UUID(idv7), userId, email, false, true, "google")
			if err != nil {
				return nil, "", ErrInternalServer
			}
			return s.provideSessionId(email)
		}
		if err1 != nil {
			return nil, "", ErrInternalServer
		}
		err = s.repository.SaveThirdPartySignInInfo(ctx, id1, userId, email, phoneNumberVerified, emailVerified, "google")
		if err != nil {
			return nil, "", ErrInternalServer
		}
		if !phoneNumberVerified {
			return s.provideSessionId(email)
		}
		return s.provideTokens(id1, role1)
	}
	if err != nil {
		return nil, "", ErrSignInWithGoogle
	}
	if emailFromDB != email {
		err = s.repository.UpdateThirdPartyEmail(context.Background(), email, userId, "apple")
		if err != nil {
			return nil, "", ErrInternalServer
		}
	}
	phoneNumberVerified, err := s.repository.FindPhoneNumberVerifiedById(id)
	if err != nil {
		return nil, "", ErrSignInWithApple
	}
	if !phoneNumberVerified {
		return s.provideSessionId(email)
	}
	return s.provideTokens(id, role)
}
