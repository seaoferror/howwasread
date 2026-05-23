package repository

import (
	"log"
	"os"

	"github.com/valkey-io/valkey-go"
)

type Repository struct {
	client valkey.Client
}

func NewRepository() *Repository {
	clientOption := valkey.ClientOption{
		Username:    os.Getenv("VALKEY_USERNAME"),
		Password:    os.Getenv("VALKEY_PASSWORD"),
		InitAddress: []string{os.Getenv("VALKEY_ADDRESS")},
	}
	client, err := valkey.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := Repository{
		client: client,
	}

	return &r
}
