package service

import (
	"backend/common/producer"
	"backend/notification/internal/repository"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	repository    *repository.Repository
	producer      *producer.Producer
	signer        *sign.URLSigner
	cloudfrontURL string
}

func NewService(r *repository.Repository, p *producer.Producer) *Service {
	pk, err := sign.LoadPEMPrivKeyFile("aws-cloudfront-private-key.pem")
	if err != nil {
		log.Panicf("fail to make cloud front private key: %v", err)
	}
	signer := sign.NewURLSigner(os.Getenv("AWS_CLOUDFRONT_KEY_ID"), pk)

	s := Service{
		repository:    r,
		producer:      p,
		signer:        signer,
		cloudfrontURL: os.Getenv("AWS_CLOUD_FRONT_URL"),
	}
	return &s
}
