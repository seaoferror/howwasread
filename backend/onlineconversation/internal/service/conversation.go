package service

import (
	"backend/common/payload"
	"backend/onlineconversation/internal/dto"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (s *Service) CreateConversation(
	ctx context.Context,
	memberId uuid.UUID,
	novel,
	shortStory,
	poem,
	play,
	film,
	by,
	rule string,
	capacity int,
	when time.Time,
	length time.Duration,
) (map[string]uuid.UUID, error) {
	conversationId, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid v7 for online conversation id", "err", err)
		return nil, err
	}
	err = s.repository.SaveConversation(
		ctx,
		memberId,
		conversationId,
		novel,
		shortStory,
		poem,
		play,
		film,
		by,
		rule,
		capacity,
		when,
		length,
	)
	if err != nil {
		return nil, err
	}
	slog.Info("success to create conversation")
	return map[string]uuid.UUID{"conversationId": conversationId}, nil
}

func (s *Service) GetConversations(ctx context.Context, memberId uuid.UUID, page int) ([]dto.ConversationFeedResponse, error) {
	resp := []dto.ConversationFeedResponse{}

	items, err := s.repository.FindConversations(ctx, page)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		var isBanned bool
		for _, bId := range item.BanIds {
			if bytes.Equal(bId.Data, memberId[:]) {
				isBanned = true
			}
		}
		if isBanned {
			continue
		}
		var isModerator bool
		var isRegistrant bool
		mIds := make([]uuid.UUID, 0, len(item.ModeratorIds))
		for _, mId := range item.ModeratorIds {
			mIds = append(mIds, uuid.UUID(mId.Data))
			if bytes.Equal(mId.Data, memberId[:]) {
				isModerator = true
			}
		}
		for _, rId := range item.RegistrantIds {
			if bytes.Equal(rId.Data, memberId[:]) {
				isRegistrant = true
				break
			}
		}
		var ongoing bool
		if time.Now().After(item.When.Add(-10 * time.Minute)) {
			ongoing = true
		}

		resp = append(resp, dto.ConversationFeedResponse{
			Id:           uuid.UUID(item.Id.Data),
			Novel:        item.Novel,
			ShortStory:   item.ShortStory,
			Poem:         item.Poem,
			Play:         item.Play,
			Film:         item.Film,
			By:           item.By,
			Rule:         item.Rule,
			Capacity:     item.Capacity,
			When:         item.When,
			Length:       item.Length.String(),
			Ongoing:      ongoing,
			IsModerator:  isModerator,
			IsRegistrant: isRegistrant,
			ModeratorIds: mIds,
		})
	}
	slog.Info("success to get conversation")
	return resp, nil
}

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

func (s *Service) PublishConversationSignal(fromId uuid.UUID, toIds [][]byte, signal []byte) error {
	value := payload.Marshal(payload.OnlineConversationSignal{
		FromId: fromId[:],
		ToIds:  toIds,
		Signal: signal,
	})
	s.producer.PushMessage("conversation-signal", nil, value)
	return nil
}

func (s *Service) GetConversation(ctx context.Context, conversationId, memberId uuid.UUID) (*dto.GetConversationResponse, error) {
	c, err := s.repository.GetConversation(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	var isModerator bool
	for _, m := range c.ModeratorIds {
		if bytes.Equal(m.Data, memberId[:]) {
			isModerator = true
			break
		}
	}

	resp := dto.GetConversationResponse{
		Id:          uuid.UUID(c.Id.Data),
		Novel:       c.Novel,
		ShortStory:  c.ShortStory,
		Poem:        c.Poem,
		Play:        c.Play,
		Film:        c.Film,
		By:          c.By,
		Rule:        c.Rule,
		When:        c.When,
		Length:      c.Length.String(),
		IsModerator: isModerator,
	}

	return &resp, nil
}

func (s *Service) BanParticipant(ctx context.Context, modId, conversationId, banId uuid.UUID) error {
	mIds, err := s.repository.FindModeratorIds(ctx, conversationId)
	if err != nil {
		return err
	}
	isMod := false
	for _, mId := range mIds {
		if bytes.Equal(mId.Data, modId[:]) {
			isMod = true
			break
		}
	}
	if !isMod {
		return errors.New("you cannot ban")
	}
	err = s.repository.AddBanId(ctx, conversationId, banId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ReportOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	ids, err := s.repository.FindReporterIdsByConversationId(ctx, conversationId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if bytes.Equal(id.Data, memberId[:]) {
			return nil
		}
	}
	if len(ids) > 5 {
		err = s.repository.DeleteOnlineConversation(ctx, conversationId)
		if err != nil {
			return err
		}
	}
	err = s.repository.AddReporterIdByConversationId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}
