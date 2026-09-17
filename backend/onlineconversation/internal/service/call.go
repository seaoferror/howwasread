package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (s *Service) GetParticipantsWithoutMe(ctx context.Context, conversationId string, memberId uuid.UUID) ([]uuid.UUID, error) {
	pidRaws, err := s.repository.FindParticipantIds(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	pids := make([]uuid.UUID, 0, len(pidRaws))
	for _, pidRaw := range pidRaws {
		pid, err1 := uuid.FromBytes([]byte(pidRaw))
		if err1 != nil {
			slog.Error("fail to parse uuid from pidRaw",
				"err", err1,
				"pidRaw", pidRaw)
			return nil, err1
		}
		if memberId == pid {
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (s *Service) AddParticipant(ctx context.Context, conversationId string, memberId uuid.UUID) error {
	err := s.repository.AddParticipantId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveParticipant(ctx context.Context, conversationId string, memberId uuid.UUID) error {
	err := s.repository.RemoveParticipantId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SetServerIP(ctx context.Context, memberId uuid.UUID, ip string) error {
	err := s.repository.SetServerIP(ctx, string(memberId[:]), ip)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveServerIP(ctx context.Context, memberId uuid.UUID) error {
	err := s.repository.RemoveServerIP(ctx, string(memberId[:]))
	if err != nil {
		return err
	}
	return nil
}
