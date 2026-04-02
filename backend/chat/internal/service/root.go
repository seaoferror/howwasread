package service

import (
	"backend/chat/internal/kafka/producer"
	"backend/chat/internal/repository"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	repository    *repository.Repository
	kafkaProducer *producer.KafkaProducer
	clientConn    *grpc.ClientConn
}

func NewService(r *repository.Repository, kp *producer.KafkaProducer) *Service {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	clientConn, err := grpc.NewClient("push-service:50051", opts...)
	if err != nil {
		slog.Error("fail to get *ClientConn",
			"err", err)
		panic(err)
	}
	s := Service{
		repository:    r,
		kafkaProducer: kp,
		clientConn:    clientConn,
	}
	return &s
}
