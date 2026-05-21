package repository

import (
	"log"
	"os"

	"github.com/redis/rueidis"
)

type Repository struct {
	client rueidis.Client
}

func NewRepository() *Repository {
	clientOption := rueidis.ClientOption{
		InitAddress: []string{os.Getenv("REDIS_ADDRESS")},
	}
	client, err := rueidis.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := Repository{
		client: client,
	}

	return &r
}
