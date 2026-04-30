package service

import (
	"backend/notification/internal/repository"
	"context"
	"crypto/ecdsa"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/option"

	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	repository       *repository.Repository
	app              *firebase.App
	teamId           string
	BundleIdentifier string
	keyId            string
	privateKey       *ecdsa.PrivateKey
}

func NewService(r *repository.Repository) *Service {
	opt := option.WithCredentialsFile("holiday2-3d1c4-firebase-adminsdk-fbsvc-7d8ad6c905.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	keyRaw, err := os.ReadFile("AuthKey_6ZYM46U7FP.p8")
	if err != nil {
		log.Panicf("Failed to read apn private key file: %v", err)
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(keyRaw)
	if err != nil {
		log.Panicf("Fail to parse apn private key form bytes: %v", err)
	}

	s := Service{
		repository:       r,
		app:              app,
		teamId:           os.Getenv("TEAM_ID"),
		BundleIdentifier: os.Getenv("BUNDLE_IDENTIFIER"),
		keyId:            os.Getenv("KEY_ID"),
		privateKey:       privateKey,
	}
	return &s
}
