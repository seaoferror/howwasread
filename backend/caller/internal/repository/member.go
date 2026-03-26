package repository

import (
	"backend/caller/internal/document"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) GetServerIP(ctx context.Context, id uuid.UUID) (string, error) {
	filter := bson.M{
		"_id": bson.Binary{Subtype: 4, Data: id[:]},
	}

	opts := options.FindOne().SetProjection(bson.M{
		"server_ip": 1,
		"_id":       0,
	})

	var d document.Member
	err := r.db.Collection("member").FindOne(ctx, filter, opts).Decode(&d)
	if err != nil {
		slog.Info("fail to find ip",
			"err", err)
		return "", err
	}
	return d.ServerIP, nil
}
