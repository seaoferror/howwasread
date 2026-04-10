package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) SetServerIP(ctx context.Context, memberId uuid.UUID, ip string) error {
	err := s.repository.SetServerIP(ctx, string(memberId[:]), ip)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveServerIP(ctx context.Context, memberId []byte, ip string) error {
	err := s.repository.RemoveServerIP(ctx, string(memberId), ip)
	if err != nil {
		return err
	}
	return nil
}
