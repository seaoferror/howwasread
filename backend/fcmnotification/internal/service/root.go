package service

import (
	"backend/fcmnotification/internal/repository"
	"context"
	"encoding/base64"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Service struct {
	repository *repository.Repository
	fcmClient  *messaging.Client
}

func NewService(r *repository.Repository) *Service {
	jsonRaw, err := base64.StdEncoding.DecodeString(
		os.Getenv("FIREBASE_ADMINSDK_FBSVC_JSON_BASE64"))
	if err != nil {
		panic(fmt.Errorf("failed to decode base64 secret: %w", err))
	}
	opt := option.WithCredentialsJSON(jsonRaw)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		panic(err)
	}

	s := Service{
		repository: r,
		fcmClient:  fcmClient,
	}
	return &s
}
