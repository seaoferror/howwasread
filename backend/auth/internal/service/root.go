package service

import (
	"backend/auth/internal/repository"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/twilio/twilio-go"
)

type Service struct {
	repository   *repository.Repository
	privateKeyAT *rsa.PrivateKey
	privateKeyRT *rsa.PrivateKey
	publicKeyRT  *rsa.PublicKey
	issuer       string
	audience     string
	twilioClient *twilio.RestClient
}

func NewService(r *repository.Repository) *Service {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	apiKey := os.Getenv("TWILIO_API_KEY")
	apiSecret := os.Getenv("TWILIO_API_SECRET")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   apiKey,
		Password:   apiSecret,
		AccountSid: accountSid,
	})

	return &Service{
		repository:   r,
		privateKeyAT: loadRSAPrivateKey("cert/authentication/private-key-at.pem"),
		privateKeyRT: loadRSAPrivateKey("cert/authentication/private-key-rt.pem"),
		publicKeyRT:  loadRSAPublicKey("cert/authentication/public-key-rt.pem"),
		issuer:       os.Getenv("ISSUER"),
		audience:     os.Getenv("BUNDLE_IDENTIFIER"),
		twilioClient: client,
	}
}

func loadRSAPrivateKey(filepath string) *rsa.PrivateKey {
	keyBytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Panicf("failed to read private key file: %v", err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		log.Panic("failed to decode PEM block from file")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Panicf("failed to parse RSA private key: %v", err)
	}

	return privateKey
}

func loadRSAPublicKey(filepath string) *rsa.PublicKey {
	keyBytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Panicf("failed to read public key file: %v", err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		log.Panic("failed to decode PEM block from file")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Panicf("failed to parse RSA public key: %v", err)
	}

	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		log.Panic("key is not an RSA public key")
	}
	return publicKey
}
