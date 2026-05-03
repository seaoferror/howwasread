package service

import (
	"backend/common/producer"
	"backend/onlineconversation/internal/repository"
)

type Service struct {
	repository *repository.Repository
	producer   *producer.Producer
}

func NewService(r *repository.Repository, p *producer.Producer) *Service {
	s := &Service{
		repository: r,
		producer:   p,
	}
	return s
}
