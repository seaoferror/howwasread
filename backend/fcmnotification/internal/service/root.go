package service

import (
	"backend/fcmnotification/internal/repository"
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Service struct {
	repository *repository.Repository
	fcmClient  *messaging.Client
}

func NewService(r *repository.Repository) *Service {
	opt := option.WithCredentialsFile("cert/firebase/firebase-adminsdk.json")
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
