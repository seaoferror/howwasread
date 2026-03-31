package service

import (
	"backend/caller/internal/repository"
	"sync"

	"google.golang.org/grpc"
)

type Service struct {
	repository  *repository.Repository
	clientConns map[string]*grpc.ClientConn
	ccsMutex    *sync.RWMutex
}

func NewService(r *repository.Repository) *Service {
	s := &Service{
		repository:  r,
		clientConns: make(map[string]*grpc.ClientConn),
		ccsMutex:    &sync.RWMutex{},
	}
	return s
}
