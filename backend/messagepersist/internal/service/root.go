package service

import (
	"backend/messagepersist/internal/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(r *repository.Repository) *Service {
	s := &Service{
		repository: r,
	}
	return s
}
