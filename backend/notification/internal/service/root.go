package service

import (
	"backend/notification/internal/repository"
	"context"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type Service struct {
	repository *repository.Repository
	app        *firebase.App
}

func NewService(r *repository.Repository) *Service {
	opt := option.WithCredentialsFile("backend/holiday2-3d1c4-firebase-adminsdk-fbsvc-7d8ad6c905.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
	s := Service{
		repository: r,
		app:        app,
	}
	return &s
}
