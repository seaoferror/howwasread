package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

// Keys signs access and refresh tokens and verifies refresh tokens
type Keys struct {
	PrivateKeyAT *rsa.PrivateKey
	PrivateKeyRT *rsa.PrivateKey
	PublicKeyRT  *rsa.PublicKey
}

// LoadKeys reads the token keys from the mounted cert files
func LoadKeys() Keys {
	return Keys{
		PrivateKeyAT: loadRSAPrivateKey("cert/authentication/private-key-at.pem"),
		PrivateKeyRT: loadRSAPrivateKey("cert/authentication/private-key-rt.pem"),
		PublicKeyRT:  loadRSAPublicKey("cert/authentication/public-key-rt.pem"),
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
