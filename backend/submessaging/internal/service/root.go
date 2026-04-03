package service

import (
	"backend/submessaging/internal/kafka/producer"
	"backend/submessaging/internal/repository"
	"log/slog"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	clientConn, err := grpc.NewClient("push-service:50051", opts...)
	if err != nil {
		slog.Error("fail to get *ClientConn",
			"err", err)
		panic(err)
	}
	s.clientConns["push-service"] = clientConn
	return s
}
