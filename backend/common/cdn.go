package common

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
)

// CDNClient hands out time limited CloudFront URLs for uploaded media
type CDNClient interface {
	SignedURL(contentType, filename string) (string, error)
}

type cdnClient struct {
	urlSigner     *sign.URLSigner
	cloudfrontURL string
}

func NewCDNClient() CDNClient {
	pk, err := sign.LoadPEMPrivKeyFile("cert/aws/aws-cloudfront-private-key.pem")
	if err != nil {
		log.Panicf("fail to make cloud front private key: %v", err)
	}
	return &cdnClient{
		urlSigner:     sign.NewURLSigner(os.Getenv("AWS_CLOUDFRONT_KEY_ID"), pk),
		cloudfrontURL: os.Getenv("AWS_CLOUDFRONT_URL"),
	}
}

// SignedURL signs <cloudfront url>/<content type>/<filename> for an hour
func (c *cdnClient) SignedURL(contentType, filename string) (string, error) {
	signedURL, err := c.urlSigner.Sign(
		fmt.Sprintf("%s/%s/%s", c.cloudfrontURL, contentType, filename),
		time.Now().Add(1*time.Hour))
	if err != nil {
		slog.Error("fail to generate signed URL", "err", err)
		return "", err
	}
	return signedURL, nil
}
