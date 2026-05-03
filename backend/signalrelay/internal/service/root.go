package service

import (
	"backend/common/producer"
	"backend/signalrelay/internal/repository"
	"sync"

	"google.golang.org/grpc"
)

type Service struct {
	repository  *repository.Repository
	clientConns map[string]*grpc.ClientConn
	ccsMutex    *sync.RWMutex
	producer    *producer.Producer
}

func NewService(r *repository.Repository, p *producer.Producer) *Service {
	s := &Service{
		repository:  r,
		clientConns: make(map[string]*grpc.ClientConn),
		ccsMutex:    &sync.RWMutex{},
		producer:    p,
	}
	return s
}
