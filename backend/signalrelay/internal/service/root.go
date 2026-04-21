package service

import (
	"backend/signalrelay/internal/kafka/producer"
	"backend/signalrelay/internal/repository"
	"sync"

	"google.golang.org/grpc"
)

type Service struct {
	repository  *repository.Repository
	clientConns map[string]*grpc.ClientConn
	ccsMutex    *sync.RWMutex
	producer    *producer.KafkaProducer
}

func NewService(r *repository.Repository, kp *producer.KafkaProducer) *Service {
	s := &Service{
		repository:  r,
		clientConns: make(map[string]*grpc.ClientConn),
		ccsMutex:    &sync.RWMutex{},
		producer:    kp,
	}
	return s
}
