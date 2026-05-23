package repository

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"os"

	"github.com/valkey-io/valkey-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
	mongoClient  *mongo.Client
	db           *mongo.Database
	valkeyClient valkey.Client
}

func NewRepository() *Repository {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		log.Panicf("fail to connect mongodb: %v", err)
	}
	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Panicf("fail to ping mongodb: %v", err)
	}
	slog.Info("success to connect mongodb")

	clientOption := valkey.ClientOption{
		Username:    os.Getenv("VALKEY_USERNAME"),
		Password:    os.Getenv("VALKEY_PASSWORD"),
		InitAddress: []string{os.Getenv("VALKEY_ADDRESS")},
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	v, err := valkey.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := &Repository{
		mongoClient:  mongoClient,
		db:           mongoClient.Database("db"),
		valkeyClient: v,
	}

	return r
}
