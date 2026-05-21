package service

import (
	"backend/chat/internal/repository"
	"backend/common/producer"
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	repository    *repository.Repository
	producer      *producer.Producer
	presignClient *s3.PresignClient
	signer        *sign.URLSigner
	bucketName    string
	cloudfrontURL string
}

func NewService(r *repository.Repository, kp *producer.Producer) *Service {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(os.Getenv("AWS_S3_REGIONS")))
	if err != nil {
		log.Panicf("fail to config for presigned client: %v", err)
	}
	presignClient := s3.NewPresignClient(s3.NewFromConfig(cfg))

	pk, err := sign.LoadPEMPrivKeyFile("aws-cloudfront-private-key.pem")
	if err != nil {
		log.Panicf("fail to make cloud front private key: %v", err)
	}
	signer := sign.NewURLSigner(os.Getenv("AWS_CLOUDFRONT_KEY_ID"), pk)

	s := Service{
		repository:    r,
		producer:      kp,
		presignClient: presignClient,
		signer:        signer,
		bucketName:    os.Getenv("AWS_S3_BUCKET_NAME"),
		cloudfrontURL: os.Getenv("AWS_CLOUD_FRONT_URL"),
	}
	return &s
}
