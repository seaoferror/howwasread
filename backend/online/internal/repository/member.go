package repository

import (
	"backend/online/internal/document"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) SaveNewMemberId(idRaw []byte) error {
	id := bson.Binary{
		Subtype: 0x04,
		Data:    idRaw,
	}
	doc := document.Member{Id: id}
	result, err := r.db.Collection("member").InsertOne(context.Background(), doc)
	if err != nil {
		slog.Error("fail to save new member id",
			"err", err,
			"id.Data", id.Data,
		)
		return err
	}
	slog.Info("success to save new member id",
		"result", result,
	)

	return nil
}

func (r *Repository) SetServerIP(ctx context.Context, memberId, ip string) error {
	result := r.redisClient.Do(ctx, r.redisClient.B().Set().Key(memberId).Value(ip).Build())
	if result.Error() != nil {
		slog.Error("fail to save member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) RemoveServerIP(ctx context.Context, memberId string) error {
	result := r.redisClient.Do(ctx, r.redisClient.B().Del().Key(memberId).Build())
	if result.Error() != nil {
		slog.Error("fail to remove member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) SetName(ctx context.Context, memberId uuid.UUID, name string) error {
	_, err := r.db.Collection("member").
		UpdateOne(ctx, bson.M{"_id": bson.Binary{Subtype: 4, Data: memberId[:]}},
			bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		slog.Error("fail to set name to member doc", "err", err)
		return err
	}
	return nil
}

func (r *Repository) FindProfile(ctx context.Context, memberId uuid.UUID) (*document.Member, error) {
	opt := options.FindOne().SetProjection(bson.M{"name": 1})

	var d document.Member
	err := r.db.Collection("member").FindOne(ctx, bson.M{"_id": bson.Binary{Subtype: 4, Data: memberId[:]}}, opt).Decode(&d)
	if err != nil {
		slog.Error("fail to find member profile")
	}
	return &d, nil
}
