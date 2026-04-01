package service

import (
	"backend/auth/internal/kafka/producer"
	"backend/auth/internal/repository"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/twilio/twilio-go"
)

type Service struct {
	repository    *repository.Repository
	kafkaProducer *producer.KafkaProducer
	secretKeyAT   []byte
	secretKeyRT   []byte
	issuer        string
	audience      string
	twilioClient  *twilio.RestClient
}

func NewService(r *repository.Repository, kp *producer.KafkaProducer) *Service {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	apiKey := os.Getenv("TWILIO_API_KEY")
	apiSecret := os.Getenv("TWILIO_API_SECRET")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   apiKey,
		Password:   apiSecret,
		AccountSid: accountSid,
	})

	return &Service{
		repository:    r,
		kafkaProducer: kp,
		secretKeyAT:   []byte(os.Getenv("SECRET_KEY_AT")),
		secretKeyRT:   []byte(os.Getenv("SECRET_KEY_RT")),
		issuer:        os.Getenv("ISSUER"),
		audience:      os.Getenv("BUNDLE_IDENTIFIER"),
		twilioClient:  client,
	}
}
