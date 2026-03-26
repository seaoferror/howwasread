package internal

import (
	"backend/auth/internal/kafka/producer"
	"backend/auth/internal/logger"
	"backend/auth/internal/network"
	"backend/auth/internal/repository"
	"backend/auth/internal/service"
)

func NewServer() {

	kp := producer.NewKafkaProducer()

	logger.SetLogger(kp)

	r := repository.NewRepository()

	s := service.NewService(r, kp)

	n := network.NewNetwork(s)

	n.Start()
}
