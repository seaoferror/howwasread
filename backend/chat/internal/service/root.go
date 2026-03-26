package service

import (
	"backend/chat/internal/kafka/producer"
	"backend/chat/internal/repository"
)

type Service struct {
	repository    *repository.Repository
	kafkaProducer *producer.KafkaProducer
}

func NewService(r *repository.Repository, kp *producer.KafkaProducer) *Service {
	s := Service{r, kp}
	return &s
}
