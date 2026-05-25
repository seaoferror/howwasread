package repository

import (
	"crypto/tls"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/valkey-io/valkey-go"
)

type Repository struct {
	client valkey.Client
}

func NewRepository() *Repository {
	clientOption := valkey.ClientOption{
		InitAddress: []string{os.Getenv("VALKEY_ADDRESS")},
	}
	if os.Getenv("PROFILE") == "production" {
		clientOption.Username = os.Getenv("VALKEY_USERNAME")
		clientOption.Password = os.Getenv("VALKEY_PASSWORD")
		clientOption.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	client, err := valkey.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to valkey: %v", err)
	}
	log.Print("success to connect valkey")

	r := Repository{client: client}

	return &r
}
