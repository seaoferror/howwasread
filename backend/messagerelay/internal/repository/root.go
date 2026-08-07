package repository

import (
	"backend/common/tlsconfig"
	"log"
	"os"

	"github.com/valkey-io/valkey-go"
)

type Repository struct {
	client valkey.Client
}

func NewRepository() *Repository {
	clientOption := valkey.ClientOption{
		InitAddress: []string{os.Getenv("VALKEY_ADDRESS")},
	}
	tlsConfig, err := tlsconfig.Create("", "", os.Getenv("VALKEY_CA_CERT_PATH"))
	if err != nil {
		log.Panicf("fail to create tls config for valkey")
	}
	tlsConfig.ServerName = os.Getenv("VALKEY_HOST")
	clientOption.TLSConfig = tlsConfig

	clientOption.Username = os.Getenv("VALKEY_USERNAME")
	clientOption.Password = os.Getenv("VALKEY_PASSWORD")

	client, err := valkey.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := Repository{
		client: client,
	}

	return &r
}
