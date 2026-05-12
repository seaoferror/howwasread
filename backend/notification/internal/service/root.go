package service

import (
	"backend/common/producer"
	"backend/notification/internal/repository"
	"bytes"
	"encoding/base64"
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
	pemRaw, err := base64.StdEncoding.DecodeString(
		os.Getenv("AWS_CLOUDFRONT_PRIVATE_KEY_PEM_BASE64"))
	if err != nil {
		log.Panicf("failed to decode base64 PEM: %v", err)
	}
	pemReader := bytes.NewReader(pemRaw)
	pk, err := sign.LoadPEMPrivKey(pemReader)
	if err != nil {
		log.Panicf("fail to make cloud front private key: %v", err)
	}
	signer := sign.NewURLSigner(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID"), pk)

	s := Service{
		repository:    r,
		producer:      p,
		signer:        signer,
		cloudfrontURL: os.Getenv("AWS_CLOUD_FRONT_URL"),
	}
	return &s
}
