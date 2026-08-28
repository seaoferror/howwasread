package service

import (
	"fmt"
	"log/slog"
	"time"
)

func (s *Service) generateSignedURL(contentType, filename string) (string, error) {
	signedURL, err := s.signer.Sign(
		fmt.Sprintf("%s/%s/%s",
			s.cloudfrontURL, contentType, filename),
		time.Now().Add(1*time.Hour))
	if err != nil {
		slog.Error("fail to generate signed URL",
			"err", err)
		return "", err
	}
	slog.Info("success to sign url", "signedURL", signedURL)
	return signedURL, nil
}
