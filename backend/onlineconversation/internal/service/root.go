package service

import (
	"backend/onlineconversation/internal/kafka/producer"
	"backend/onlineconversation/internal/repository"
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
