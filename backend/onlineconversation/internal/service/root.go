package service

import (
	"backend/common/producer"
	"backend/onlineconversation/internal/repository"
	"net"
	"os"
	"time"
)

type Service struct {
	repository *repository.Repository
	producer   *producer.Producer
	turnSecret string
	publicIP   net.IP
}

func NewService(r *repository.Repository, p *producer.Producer) *Service {
	s := &Service{
		repository: r,
		producer:   p,
		turnSecret: os.Getenv("TURN_SECRET"),
	}

	go s.monitorAndReflectIPChange(10 * time.Second)
	return s
}
