package service

import (
	"backend/common/payload"
	"backend/onlineconversation/internal/dto"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
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
	s.producer.PushMessage("search",
		nil,
		payload.Marshal(dto.OnlineConversationDocument{
			Id:         conversationId,
			Novel:      novel,
			ShortStory: shortStory,
			Poem:       poem,
			Play:       play,
			Film:       film,
			WrittenBy:  by,
			Time:       when,
		}),
		[]sarama.RecordHeader{
			{Key: []byte("type"), Value: []byte("onlineconversation")},
		},
	)
	return map[string]uuid.UUID{"conversationId": conversationId}, nil
}

func (s *Service) GetConversations(ctx context.Context, page int, t time.Time) ([]dto.OnlineConversationFeedResponse, error) {
	resp := []dto.OnlineConversationFeedResponse{}

	items, err := s.repository.FindConversations(ctx, page, t)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		resp = append(resp, dto.OnlineConversationFeedResponse{
			Id:         uuid.UUID(item.Id.Data),
			Novel:      item.Novel,
			ShortStory: item.ShortStory,
			Poem:       item.Poem,
			Play:       item.Play,
			Film:       item.Film,
			By:         item.By,
			When:       item.When,
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
	s.producer.PushMessage("conversation-signal", nil, value, nil)
	return nil
}

func (s *Service) GetConversationDetail(ctx context.Context, conversationId, memberId uuid.UUID) (*dto.OnlineConversationDetailResponse, error) {
	c, err := s.repository.FindConversation(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	var isRegistrant bool
	for _, r := range c.RegistrantIds {
		if bytes.Equal(r.Data, memberId[:]) {
			isRegistrant = true
			break
		}
	}
	modIds := make([]uuid.UUID, 0, len(c.ModeratorIds))
	for _, m := range c.ModeratorIds {
		modIds = append(modIds, uuid.UUID(m.Data))
	}

	canEnter := true
	if time.Now().UTC().Before(c.When.Add(-15 * time.Minute)) {
		canEnter = false
	}
	if time.Now().UTC().Before(c.When.Add(10*time.Minute)) && !isRegistrant {
		canEnter = false
	}

	resp := dto.OnlineConversationDetailResponse{
		Id:           uuid.UUID(c.Id.Data),
		Novel:        c.Novel,
		ShortStory:   c.ShortStory,
		Poem:         c.Poem,
		Play:         c.Play,
		Film:         c.Film,
		By:           c.By,
		Rule:         c.Rule,
		Capacity:     c.Capacity,
		When:         c.When,
		Length:       c.Length.String(),
		CanEnter:     canEnter,
		ModeratorIds: modIds,
		IsRegistrant: isRegistrant,
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

func (s *Service) RegisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	capacity, err := s.repository.FindCapacity(ctx, conversationId)
	if err != nil {
		return err
	}
	err = s.repository.AddRegistrantId(ctx, memberId, conversationId, capacity)
	if err != nil {
		return err
	}
	//TODO: publish message to notification with header
	return nil
}

func (s *Service) DeregisterOnlineConversation(ctx context.Context, memberId, conversationId uuid.UUID) error {
	err := s.repository.RemoveRegistrantId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}
