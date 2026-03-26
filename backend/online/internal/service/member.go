package service

import (
	"backend/online/internal/dto"
	"context"
	"errors"
	"log/slog"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

func (s *Service) SaveNewMemberId(idRaw []byte) error {
	err := s.repository.SaveNewMemberId(idRaw)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SetServerIP(ctx context.Context, memberId uuid.UUID, ip string) error {
	err := s.repository.SetServerIP(ctx, memberId, ip)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveServerIP(ctx context.Context, memberId uuid.UUID) error {
	err := s.repository.RemoveServerIP(ctx, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SetName(ctx context.Context, memberId uuid.UUID, name string) error {
	sanitizedName := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, name)
	if len(sanitizedName) == 0 {
		slog.Info("incorrect name")
		return errors.New("incorrect name")
	}
	err := s.repository.SetName(ctx, memberId, sanitizedName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetMyProfile(ctx context.Context, id uuid.UUID) (*dto.GetMyProfileResponse, error) {
	profile, err := s.repository.FindProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	res := dto.GetMyProfileResponse{Id: id.String(), Name: profile.Name}

	return &res, nil
}
