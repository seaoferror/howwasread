package client

import (
	"context"
	"errors"
	"log"
	"os"

	"google.golang.org/api/idtoken"
)

// GoogleAuthClient verifies Google sign in ID tokens
type GoogleAuthClient interface {
	// Email returns the email of a valid ID token issued for this app
	Email(ctx context.Context, idToken string) (string, error)
}

type googleAuthClient struct {
	validator *idtoken.Validator
	audience  string
}

func NewGoogleAuthClient() GoogleAuthClient {
	validator, err := idtoken.NewValidator(context.Background())
	if err != nil {
		log.Fatalf("Failed to create google id token validator: %v", err)
	}
	return &googleAuthClient{validator: validator, audience: os.Getenv("GOOGLE_SIGN_IN_WEB_CLIENT_ID")}
}

func (g *googleAuthClient) Email(ctx context.Context, idToken string) (string, error) {
	payload, err := g.validator.Validate(ctx, idToken, g.audience)
	if err != nil {
		return "", err
	}
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", errors.New("no email in google id token")
	}
	return email, nil
}
