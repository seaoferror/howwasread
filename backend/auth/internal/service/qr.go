package service

import (
	"crypto/rand"
	"encoding/base32"
	"log"
	"net/url"

	"rsc.io/qr"
)

func (s *Service) generateUserQR(email, secret string) ([]byte, error) {
	u, err := url.Parse("otpauth://totp")
	if err != nil {
		log.Printf("fail to parse totp base url: %v", err)
		return nil, err
	}
	u.Path += "/" + url.PathEscape(email)

	p := url.Values{}
	p.Add("secret", secret)
	p.Add("issuer", s.issuer)
	u.RawQuery = p.Encode()

	c, err := qr.Encode(u.String(), qr.Q)
	if err != nil {
		log.Printf("fail to create qr from url: %v", err)
		return nil, err
	}
	return c.PNG(), err
}

func generateQRSecret() (string, error) {
	bytes := make([]byte, 10) // 10 bytes, so it is 16 base32 chars
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("fail to generate random byte: %v", err)
		return "", err
	}
	secret := base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(bytes)
	return secret, nil
}
