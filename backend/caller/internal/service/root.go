package service

import (
	"backend/caller/internal/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(r *repository.Repository) *Service {
	s := &Service{r}
	return s
}
