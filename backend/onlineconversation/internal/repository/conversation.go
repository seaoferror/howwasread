package repository

import (
	"backend/onlineconversation/internal/document"
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const Limit = 10

func (r *Repository) SaveConversation(ctx context.Context, memberId, conversationId uuid.UUID, novel, shortStory, poem, play, film, writtenBy, rule string, capacity int, t time.Time, length time.Duration) error {
	newConversation := document.Conversation{
		Id:         bson.Binary{4, conversationId[:]},
		Novel:      novel,
		ShortStory: shortStory,
		Poem:       poem,
		Play:       play,
		Film:       film,
		WrittenBy:  writtenBy,
		Rule:       rule,
		Capacity:   capacity,
		Time:       t,
		Length:     length,
		ModeratorIds: []bson.Binary{
			{4, memberId[:]},
		},
		RegistrantIds: []bson.Binary{
			{4, memberId[:]},
		},
		BanIds:      []bson.Binary{},
		ReporterIds: []bson.Binary{},
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

func (r *Repository) FindConversations(ctx context.Context, page int, t time.Time) ([]document.Conversation, error) {
	filter := bson.M{
		"time":    bson.M{"$gt": t.Add(-9 * time.Hour)},
		"expired": false,
	}

	opts := options.Find().
		SetSort(bson.M{"time": 1}).
		SetLimit(Limit).
		SetSkip(int64((page - 1) * 5)).
		SetProjection(bson.M{
			"rule":           0,
			"capacity":       0,
			"length":         0,
			"moderator_ids":  0,
			"registrant_ids": 0,
			"ban_ids":        0,
			"reporter_ids":   0,
		})

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

func (r *Repository) FindParticipantIds(ctx context.Context, conversationId string) ([]string, error) {
	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Smembers().Key(conversationId).Build())
	if result.Error() != nil {
		slog.Error("fail to get member ip", "err", result.Error())
		return nil, result.Error()
	}
	value, err := result.AsStrSlice()
	if err != nil {
		slog.Error("fail to get ips value string slice", "err", err)
	}
	return value, nil
}

func (r *Repository) AddParticipantId(ctx context.Context, conversationId string, memberId uuid.UUID) error {
	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Sadd().Key(conversationId).Member(string(memberId[:])).Build())
	if result.Error() != nil {
		slog.Error("fail to add participant member id to conversation",
			"err", result.Error(),
			"conversationId", conversationId,
			"memberId", memberId)
		return result.Error()
	}
	return nil
}

func (r *Repository) RemoveParticipantId(ctx context.Context, conversationId string, memberId uuid.UUID) error {
	result := r.valkeyClient.Do(ctx, r.valkeyClient.B().Srem().Key(conversationId).Member(string(memberId[:])).Build())
	if result.Error() != nil {
		slog.Error("fail to remove participant member id to conversation",
			"err", result.Error(),
			"conversationId", conversationId,
			"memberId", memberId)
		return result.Error()
	}
	return nil
}

func (r *Repository) FindConversation(ctx context.Context, conversationId uuid.UUID) (*document.Conversation, error) {
	opts := options.FindOne().SetProjection(bson.M{"reporter_ids": 0})

	var d document.Conversation
	err := r.db.Collection("conversation").
		FindOne(ctx, bson.M{"_id": bson.Binary{4, conversationId[:]}}, opts).
		Decode(&d)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	return &d, nil
}

func (r *Repository) FindModeratorIds(ctx context.Context, conversationId uuid.UUID) ([]bson.Binary, error) {
	opt := options.FindOne().SetProjection(bson.M{"moderator_ids": 1})

	var d document.Conversation
	err := r.db.Collection("conversation").
		FindOne(ctx, bson.M{"_id": bson.Binary{4, conversationId[:]}}, opt).Decode(&d)
	if err != nil {
		slog.Error("fail to find mod ids", "err", err)
		return nil, err
	}
	return d.ModeratorIds, nil
}

func (r *Repository) AddBanId(ctx context.Context, conversationId uuid.UUID, banId uuid.UUID) error {
	_, err := r.db.Collection("conversation").
		UpdateOne(ctx, bson.M{"_id": bson.Binary{4, conversationId[:]}},
			bson.M{"$push": bson.M{"ban_ids": bson.Binary{4, banId[:]}}})
	if err != nil {
		slog.Error("fail to add participant member id to conversation",
			"err", err,
			"conversationId", conversationId,
			"memberId", banId.String())
		return err
	}
	return nil
}

func (r *Repository) FindReporterIdsByConversationId(ctx context.Context, conversationId uuid.UUID) ([]bson.Binary, error) {
	opt := options.FindOne().SetProjection(bson.M{"reporter_ids": 1})

	var d document.Conversation
	err := r.db.Collection("conversation").
		FindOne(ctx, bson.M{"_id": bson.Binary{4, conversationId[:]}}, opt).Decode(&d)
	if err != nil {
		slog.Error("fail to find reporter ids", "err", err)
		return nil, err
	}
	return d.ReporterIds, nil
}

func (r *Repository) DeleteOnlineConversation(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Collection("conversation").DeleteOne(ctx, bson.M{"_id": bson.Binary{4, id[:]}})
	if err != nil {
		slog.Error("fail to delete conversation",
			"err", err)
		return err
	}
	return nil
}

func (r *Repository) AddReporterIdByConversationId(ctx context.Context, conversationId, memberId uuid.UUID) error {
	filter := bson.M{"_id": bson.Binary{4, conversationId[:]}}
	update := bson.M{
		"$addToSet": bson.M{
			"reporter_ids": bson.Binary{
				Subtype: 4,
				Data:    memberId[:],
			},
		},
	}
	_, err := r.db.Collection("conversation").UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("fail to update reporter ids for conversation",
			"conversationId", conversationId,
			"memberId", memberId,
			"err", err,
		)
		return err
	}
	return nil
}

func (r *Repository) FindRegistrantIdsByConversationId(ctx context.Context, conversationId uuid.UUID) ([]bson.Binary, error) {
	opt := options.FindOne().SetProjection(bson.M{"reporter_ids": 1})

	var d document.Conversation
	err := r.db.Collection("conversation").
		FindOne(ctx, bson.M{"_id": bson.Binary{4, conversationId[:]}}, opt).Decode(&d)
	if err != nil {
		slog.Error("fail to find reporter ids", "err", err)
		return nil, err
	}
	return d.ReporterIds, nil
}

func (r *Repository) FindCapacity(ctx context.Context, id uuid.UUID) (int, error) {
	opt := options.FindOne().SetProjection(bson.M{"capacity": 1})
	var d document.Conversation
	err := r.db.Collection("conversation").
		FindOne(ctx, bson.M{"_id": bson.Binary{4, id[:]}}, opt).Decode(&d)
	if err != nil {
		slog.Error("fail to find capacity", "err", err)
		return 0, err
	}
	return d.Capacity, nil
}

var ErrMaxRegistrantReached = errors.New("already fully registered")

func (r *Repository) AddRegistrantId(ctx context.Context, conversationId, memberId uuid.UUID, capacity int) error {
	filter := bson.M{
		"_id": bson.Binary{Subtype: 4, Data: conversationId[:]},
		"registrant_ids." + strconv.Itoa(capacity-1): bson.M{"$exists": false},
	}

	update := bson.M{
		"$addToSet": bson.M{
			"registrant_ids": bson.Binary{
				Subtype: 4,
				Data:    memberId[:],
			},
		},
	}

	res, err := r.db.Collection("conversation").UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("fail to update reporter ids for conversation",
			"conversationId", conversationId,
			"memberId", memberId,
			"err", err,
		)
		return err
	}

	if res.MatchedCount == 0 {
		return ErrMaxRegistrantReached
	}

	return nil
}

func (r *Repository) RemoveRegistrantId(ctx context.Context, conversationId, memberId uuid.UUID) error {
	filter := bson.M{"_id": bson.Binary{4, conversationId[:]}}
	update := bson.M{
		"$pull": bson.M{
			"registrant_ids": bson.Binary{
				Subtype: 4,
				Data:    memberId[:],
			},
		},
	}
	_, err := r.db.Collection("conversation").UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error("fail to remove registrant id",
			"conversationId", conversationId,
			"memberId", memberId,
			"err", err)
		return err
	}
	return nil
}
