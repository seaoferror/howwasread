package repository

import (
	"backend/online/internal/document"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const Limit = 10

func (r *Repository) SaveConversation(ctx context.Context, memberId uuid.UUID, conversationId bson.ObjectID, novel, shortStory, poem, play, film, by, rule string, capacity int, when time.Time, length time.Duration) error {
	newConversation := document.Conversation{
		Id:         conversationId,
		Novel:      novel,
		ShortStory: shortStory,
		Poem:       poem,
		Play:       play,
		Film:       film,
		By:         by,
		Rule:       rule,
		Capacity:   capacity,
		When:       when,
		Length:     length,
		ModeratorIds: []bson.Binary{
			{4, memberId[:]},
		},
		ParticipantIds: []bson.Binary{},
		BanIds:         []bson.Binary{},
	}
	filter := bson.M{"_id": bson.Binary{Subtype: 4, Data: memberId[:]}}
	update := bson.M{
		"$push": bson.M{
			"m_c_ids": conversationId,
		},
	}

	session, err := r.mongoClient.StartSession()
	if err != nil {
		slog.Error("fail to start transaction for insert new conversation",
			"err", err)
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		_, err = r.db.Collection("conversation").InsertOne(ctx, newConversation)
		if err != nil {
			slog.Error("fail to insert new conversation",
				"err", err)
			return nil, err
		}

		_, err = r.db.Collection("member").UpdateOne(ctx, filter, update)
		if err != nil {
			slog.Error("fail to insert new conversation id to moderator conversation member array",
				"err", err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		slog.Error("fail to transaction for saving new conversation",
			"err", err,
		)
		return err
	}

	return nil
}

func (r *Repository) FindConversations(ctx context.Context, page int) ([]document.Conversation, error) {
	filter := bson.M{
		"when":    bson.M{"$gt": time.Now().Add(-9 * time.Hour)},
		"expired": false,
	}

	opts := options.Find().
		SetSort(bson.M{"when": 1}).
		SetLimit(Limit).
		SetSkip(int64((page - 1) * 5))

	c, err := r.db.Collection("conversation").Find(ctx, filter, opts)
	if err != nil {
		slog.Error("fail to find next conversations page",
			"err", err)
		return nil, err
	}

	items := make([]document.Conversation, 0, Limit)
	err = c.All(ctx, &items)
	if err != nil {
		return nil, err
	}
	err = c.Close(ctx)
	if err != nil {
		slog.Error("fail to close *Cursor",
			"err", err)
	}
	return items, nil
}

func (r *Repository) FindParticipantIds(ctx context.Context, conversationId bson.ObjectID) ([]bson.Binary, error) {
	opts := options.FindOne().SetProjection(bson.M{"p_ids": 1})

	var d document.Conversation
	err := r.db.Collection("conversation").FindOne(ctx, bson.M{"_id": conversationId}, opts).Decode(&d)
	if err != nil {
		slog.Error("fail to find p", "err", err)
	}
	return d.ParticipantIds, nil
}

func (r *Repository) AddParticipantId(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) error {
	_, err := r.db.Collection("conversation").
		UpdateOne(ctx, bson.M{"_id": conversationId},
			bson.M{"$push": bson.M{"p_ids": bson.Binary{4, memberId[:]}}})
	if err != nil {
		slog.Error("fail to add participant member id to conversation",
			"err", err,
			"conversationId", conversationId,
			"memberId", memberId)
		return err
	}
	return nil
}

func (r *Repository) RemoveParticipantId(ctx context.Context, conversationId bson.ObjectID, memberId uuid.UUID) error {
	_, err := r.db.Collection("conversation").
		UpdateOne(ctx, bson.M{"_id": conversationId},
			bson.M{"$pull": bson.M{"p_ids": bson.Binary{4, memberId[:]}}})
	if err != nil {
		slog.Error("fail to remove participant member id to conversation",
			"err", err,
			"conversationId", conversationId,
			"memberId", memberId)
		return err
	}
	return nil
}

func (r *Repository) GetConversation(ctx context.Context, conversationId bson.ObjectID) (*document.Conversation, error) {
	opts := options.FindOne().SetProjection(bson.M{"p_ids": 0, "r_ids": 0})

	var d document.Conversation
	err := r.db.Collection("conversation").FindOne(ctx, bson.M{"_id": conversationId}, opts).Decode(&d)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	return &d, nil
}

func (r *Repository) FindModeratorIds(ctx context.Context, conversationId bson.ObjectID) ([]bson.Binary, error) {
	opt := options.FindOne().SetProjection(bson.M{"m_ids": 1})

	var d document.Conversation
	err := r.db.Collection("conversation").
		FindOne(ctx, bson.M{"_id": conversationId}, opt).Decode(&d)
	if err != nil {
		slog.Error("fail to find mod ids", "err", err)
		return nil, err
	}
	return d.ModeratorIds, nil
}

func (r *Repository) AddBanId(ctx context.Context, conversationId bson.ObjectID, banId uuid.UUID) error {
	_, err := r.db.Collection("conversation").
		UpdateOne(ctx, bson.M{"_id": conversationId},
			bson.M{"$push": bson.M{"b_ids": bson.Binary{4, banId[:]}}})
	if err != nil {
		slog.Error("fail to add participant member id to conversation",
			"err", err,
			"conversationId", conversationId,
			"memberId", banId.String())
		return err
	}
	return nil
}
