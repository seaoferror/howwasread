package service

import (
	"backend/messagepreprocess/internal/kafka/producer"
	"backend/messagepreprocess/internal/repository"
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
