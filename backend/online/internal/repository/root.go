package repository

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/redis/rueidis"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
	mongoClient *mongo.Client
	db          *mongo.Database
	redisClient rueidis.Client
}

func NewRepository() *Repository {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		log.Panicf("fail to connect mongodb: %v", err)
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Panicf("fail to ping mongodb: %v", err)
	}
	slog.Info("success to connect mongodb")

	clientOption := rueidis.ClientOption{
		InitAddress: []string{os.Getenv("REDIS_URL")},
	}
	redisClient, err := rueidis.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := &Repository{
		mongoClient: mongoClient,
		db:          mongoClient.Database("db"),
		redisClient: redisClient,
	}

	return r
}
