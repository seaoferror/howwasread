package repository

import (
	"context"
	"log"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewRepository() *Repository {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Panicf("fail to connect mongodb: %v", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Panicf("fail to ping mongodb: %v", err)
	}
	slog.Info("success to connect mongodb")
	r := &Repository{client, client.Database("db")}
	return r
}
