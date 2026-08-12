package service

import (
	"backend/common/producer"
	"backend/onlineconversation/internal/repository"
	"os"
)

type Service struct {
	repository *repository.Repository
	producer   *producer.Producer
	turnSecret string
	turnHost   string
}

func NewService(r *repository.Repository, p *producer.Producer) *Service {
	s := &Service{
		repository: r,
		producer:   p,
		turnSecret: os.Getenv("TURN_SECRET"),
		turnHost:   os.Getenv("TURN_HOST"),
	}

	return s
}
