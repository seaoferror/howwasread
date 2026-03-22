package service

import (
	"backend/online/server/kafka/producer"
	"backend/online/server/repository"

	"github.com/google/uuid"
)

type Service struct {
	repository    *repository.Repository
	kafkaProducer *producer.KafkaProducer
}

func NewService(r *repository.Repository, kp *producer.KafkaProducer) *Service {
	s := &Service{
		repository:    r,
		kafkaProducer: kp,
	}
	return s
}
