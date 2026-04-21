package service

import (
	"backend/messageprocess/internal/kafka/producer"
	"backend/messageprocess/internal/repository"
)

type Service struct {
	repository *repository.Repository
	producer   *producer.Producer
}

func NewService(r *repository.Repository, kp *producer.Producer) *Service {
	s := &Service{
		repository: r,
		producer:   kp,
	}
	return s
}
