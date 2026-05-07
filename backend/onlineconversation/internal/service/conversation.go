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
	"go.mongodb.org/mongo-driver/v2/bson"
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
) (map[string]string, error) {
	conversationId := bson.NewObjectID()
	err := s.repository.SaveConversation(
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
	return map[string]string{"conversationId": conversationId.Hex()}, nil
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
		for _, mId := range item.ModeratorIds {
			if bytes.Equal(mId.Data, memberId[:]) {
				isModerator = true
				break
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
			Id:           item.Id.Hex(),
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
		})
	}
	slog.Info("success to get conversation")
	return resp, nil
}

func (s *Service) GetParticipants(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) (pids []uuid.UUID, err error) {
	pidsRaw, err := s.repository.FindParticipantIds(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	for _, pidRaw := range pidsRaw {
		pid, err := uuid.FromBytes(pidRaw.Data)
		if err != nil {
			slog.Error("fail to parse uuid from pidRaw",
				"err", err,
				"pidRaw.Data", pidRaw.Data)
			return nil, err
		}
		if memberId == pid {
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (s *Service) AddParticipant(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) error {
	err := s.repository.AddParticipantId(ctx, conversationId, memberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveParticipant(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) error {
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
	err := s.producer.PushMessage("conversation.signal", nil, value)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetConversation(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) (*dto.GetConversationResponse, error) {
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
		Id:          c.Id.Hex(),
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

func (s *Service) BanParticipant(ctx context.Context, modId uuid.UUID, conversationId bson.ObjectID, banId uuid.UUID) error {
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
