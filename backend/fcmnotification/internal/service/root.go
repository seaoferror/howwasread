package service

import (
	"backend/common/producer"
	"backend/fcmnotification/internal/repository"
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Service struct {
	producer   *producer.Producer
	repository *repository.Repository
	fcmClient  *messaging.Client
}

func NewService(r *repository.Repository, p *producer.Producer) *Service {
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
		producer:   p,
		repository: r,
		fcmClient:  fcmClient,
	}
	return &s
}
