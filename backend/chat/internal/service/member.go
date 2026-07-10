package service

import (
	"backend/chat/internal/dto"
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
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

func (s *Service) CheckBlock(ctx context.Context, blockerId uuid.UUID, blockedId uuid.UUID) (map[string]bool, error) {
	w, err := s.repository.DidBlock(ctx, gocql.UUID(blockerId), gocql.UUID(blockedId))
	if err != nil {
		return nil, err
	}
	return map[string]bool{"didBlock": w}, nil
}

func (s *Service) GetChatParticipants(ctx context.Context, roomId uuid.UUID) ([]dto.GetProfileResponse, error) {
	ps, err := s.repository.FindChatParticipantIds(ctx, gocql.UUID(roomId))
	if err != nil {
		return nil, err
	}
	res := make([]dto.GetProfileResponse, 0)
	for _, p := range ps {
		if p == (gocql.UUID{}) {
			continue
		}
		res = append(res, dto.GetProfileResponse{Id: uuid.UUID(p)})
	}
	return res, nil
}

func (s *Service) ReportUser(ctx context.Context, reporterId, reportedId uuid.UUID) error {
	err := s.repository.AddReporterIdByReportedId(ctx, gocql.UUID(reporterId), gocql.UUID(reportedId))
	if err != nil {
		return err
	}
	rc, err := s.repository.FindReportCountById(ctx, gocql.UUID(reportedId))
	if err != nil {
		return err
	}
	if rc > 5 {
		email, phoneNumber, err1 := s.repository.FindEmailAndPhoneNumberById(ctx, gocql.UUID(reportedId))
		if err1 != nil {
			return err1
		}
		err = s.repository.DeleteAccount(ctx, gocql.UUID(reportedId), email, phoneNumber)
		if err != nil {
			return err
		}
		err = s.repository.BanPhoneNumber(ctx, phoneNumber)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) BlockConversation(ctx context.Context, memberId uuid.UUID, conversationId string) error {
	err := s.repository.AddBlockedConversation(ctx, gocql.UUID(memberId), conversationId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetBlockedConversations(ctx context.Context, memberId uuid.UUID) ([]dto.BlockReport, error) {
	ids, err := s.repository.FindBlockedConversations(ctx, gocql.UUID(memberId))
	if err != nil {
		return nil, err
	}
	res := make([]dto.BlockReport, 0, len(ids))
	for _, id := range ids {
		res = append(res, dto.BlockReport{Id: id})
	}
	return res, nil
}
