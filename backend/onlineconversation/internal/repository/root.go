package repository

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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

	certPath := "cert/mongodb/mongodb-cert.pem"
	certs, err := tls.LoadX509KeyPair(certPath, certPath)
	if err != nil {
		log.Panicf("fail to load mongodb client certificate: %v", err)
	}
	opts := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetServerAPIOptions(serverAPI).
		SetAuth(options.Credential{
			AuthMechanism: "MONGODB-X509",
		}).
		SetTLSConfig(&tls.Config{
			Certificates: []tls.Certificate{certs},
		})
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
		InitAddress: []string{os.Getenv("VALKEY_ADDRESS")},
	}
	if os.Getenv("PROFILE") == "production" {
		caCertPEM, err1 := os.ReadFile("/cert/valkey/ca.crt")
		if err1 != nil {
			log.Fatalf("Failed to read CA certificate: %v", err1)
		}
		rootCAs := x509.NewCertPool()
		ok := rootCAs.AppendCertsFromPEM(caCertPEM)
		if !ok {
			log.Fatalf("Failed to parse root certificate")
		}
		clientOption.Username = os.Getenv("VALKEY_USERNAME")
		clientOption.Password = os.Getenv("VALKEY_PASSWORD")
		clientOption.TLSConfig = &tls.Config{
			RootCAs:            rootCAs,
			InsecureSkipVerify: false,
		}
	}
	v, err := valkey.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := &Repository{
		mongoClient:  mongoClient,
		db:           mongoClient.Database(os.Getenv("PROFILE")),
		valkeyClient: v,
	}

	return r
}
