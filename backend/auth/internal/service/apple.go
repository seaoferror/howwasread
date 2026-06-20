package service

import (
	"backend/auth/internal/constant"
	"backend/auth/internal/dto"
	"errors"
	"log/slog"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) SignInWithApple(identityToken string, isFirstSignIn bool) (*dto.SignInWithAppleResponse, string, error) {
	idt, err := jwt.Parse(identityToken, s.appleKeyFunc)
	if err != nil {
		slog.Info("fail to parse identityToken with apple JWKs", "err", err)
		return nil, "", ErrSignInWithApple
	}
	if !idt.Valid {
		return nil, "", ErrSignInWithApple
	}
	issFromClaims, err := idt.Claims.GetIssuer()
	if err != nil {
		slog.Info("fail to get issuer",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, "", ErrSignInWithApple
	}
	if issFromClaims != constant.AppleIssuerUrl {
		slog.Info("not expected apple issuer",
			"iss", issFromClaims,
		)
		return nil, "", ErrSignInWithApple
	}
	audsFromClaims, err := idt.Claims.GetAudience()
	if err != nil {
		slog.Info("fail to get audience",
			"err", err,
			"claims", idt.Claims,
		)
		return nil, "", ErrSignInWithApple
	}
	if len(audsFromClaims) == 0 || audsFromClaims[0] != s.audience {
		slog.Info("no audience or not expected audience")
		return nil, "", ErrSignInWithApple
	}
	exp, err := idt.Claims.GetExpirationTime()
	if err != nil {
		slog.Info("fail to get expiration time")
		return nil, "", ErrSignInWithApple
	}
	if exp.Unix() < time.Now().Unix() {
		slog.Info("stale token")
		return nil, "", ErrSignInWithApple
	}
	nonce, ok := idt.Claims.(jwt.MapClaims)["nonce"].(string)
	if !ok {
		slog.Info("no nonce in claims")
		return nil, "", ErrSignInWithApple
	}
	exist, err := s.repository.CheckNonce(nonce)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	if exist {
		return nil, "", ErrSignInWithApple
	}
	err = s.repository.SaveNonce(nonce)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	user, err := idt.Claims.GetSubject()
	if err != nil {
		slog.Info("fail to get subject",
			"err", err,
		)
	}
	email, ok := idt.Claims.(jwt.MapClaims)["email"].(string)
	if !ok {
		slog.Info("no email in claims")
		return nil, "", ErrSignInWithApple
	}
	if isFirstSignIn {
		emailVerified, phoneNumberVerified, id, _, role, err1 := s.repository.FindLoginInfoByEmail(email)
		if errors.Is(err1, gocql.ErrNotFound) {
			err1 = nil
			idv7, err2 := uuid.NewV7()
			if err2 != nil {
				slog.Error("fail to create uuid v7 for apple sign in user")
				return nil, "", ErrInternalServer
			}
			id = gocql.UUID(idv7)
			err = s.repository.SaveAppleSignInInfo(id, user, email, false, true)
			if err != nil {
				return nil, "", ErrInternalServer
			}
			sessionId, err2 := gocql.RandomUUID()
			if err2 != nil {
				slog.Error("fail to generate random uuid for session")
				return nil, "", ErrInternalServer
			}
			err = s.repository.SaveEmailBySessionId(sessionId, email)
			if err != nil {
				return nil, "", ErrInternalServer
			}
			resp := dto.SignInWithAppleResponse{
				SessionId: uuid.UUID(sessionId),
			}
			return &resp, "", nil
		}
		if err1 != nil {
			return nil, "", ErrSignInWithApple
		}
		err = s.repository.SaveAppleSignInInfo(id, user, email, phoneNumberVerified, emailVerified)
		if err != nil {
			return nil, "", ErrInternalServer
		}
		if !phoneNumberVerified {
			sessionId := uuid.New()
			err = s.repository.SaveEmailBySessionId(gocql.UUID(sessionId), email)
			if err != nil {
				return nil, "", ErrSignInWithApple
			}
			resp := dto.SignInWithAppleResponse{
				SessionId: sessionId,
			}
			return &resp, "", nil
		}
		jti, err1 := gocql.RandomUUID()
		if err1 != nil {
			slog.Error("fail to create random uuid for jti")
			return nil, "", ErrSignInWithApple
		}

		at, rt, err1 := s.createLoginTokens(id.String(), jti.String(), role)
		if err1 != nil {
			return nil, "", ErrSignInWithApple
		}

		err = s.repository.SaveRefreshTokenJTIById(id, jti)
		if err != nil {
			return nil, "", ErrSignInWithApple
		}
		resp := dto.SignInWithAppleResponse{
			AccessToken: at,
		}
		return &resp, rt, nil
	}

	id, emailFromDB, role, err := s.repository.FindAppleSignInInfoByUser(user)
	if err != nil {
		return nil, "", ErrSignInWithApple
	}
	//this additional fetching can be removed to improve speed little bit
	//by adding few lines, but the advantage is also small currently and
	//make link phone number process more complicate
	phoneNumberVerified, err := s.repository.FindPhoneNumberVerifiedById(id)
	if err != nil {
		return nil, "", ErrSignInWithApple
	}
	if !phoneNumberVerified {
		sessionId := uuid.New()
		err = s.repository.SaveEmailBySessionId(gocql.UUID(sessionId), emailFromDB)
		if err != nil {
			return nil, "", ErrSignInWithApple
		}
		resp := dto.SignInWithAppleResponse{
			SessionId: sessionId,
		}
		return &resp, "", nil
	}
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
	resp := dto.SignInWithAppleResponse{
		AccessToken: at,
	}
	return &resp, rt, nil
}
