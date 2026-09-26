package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// StorageClient lets clients upload media straight to S3
type StorageClient interface {
	// PresignUpload returns the form url and fields of a presigned POST for <contentType>/<filename>
	PresignUpload(ctx context.Context, contentType, filename string) (url string, fields map[string]string, err error)
}

type storageClient struct {
	presignClient *s3.PresignClient
	bucketName    string
}

func NewStorageClient() StorageClient {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(os.Getenv("AWS_S3_REGION")))
	if err != nil {
		log.Panicf("fail to config for presigned client: %v", err)
	}
	return &storageClient{
		presignClient: s3.NewPresignClient(s3.NewFromConfig(cfg)),
		bucketName:    os.Getenv("AWS_S3_BUCKET_NAME"),
	}
}

func (s *storageClient) PresignUpload(ctx context.Context, contentType, filename string) (string, map[string]string, error) {
	p, err := s.presignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", contentType, filename)),
	}, func(opts *s3.PresignPostOptions) {
		opts.Expires = 1 * time.Hour
		opts.Conditions = []any{
			[]any{"content-length-range", 1, 1024 * 1024 * 1024},
			[]any{"starts-with", "$Content-Type", contentType},
		}
	})
	if err != nil {
		return "", nil, err
	}
	return p.URL, p.Values, nil
}
