package service

import (
	"backend/chat/internal/dto"
	"backend/common/payload"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func (s *Service) SendLike(ctx context.Context, fromId, toId uuid.UUID) error {
	messageId, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to create uuid v7 for saving payload", "err", err)
		return err
	}
	err = s.repository.SaveLike(ctx, gocql.UUID(messageId), gocql.UUID(fromId), gocql.UUID(toId))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetRecentMessages(ctx context.Context, id, cursor uuid.UUID) (res []dto.MessagingResponse, err error) {
	result, err := s.repository.FindRecentMessagesByToId(ctx, gocql.UUID(id), gocql.UUID(cursor))
	if err != nil {
		return nil, err
	}
	for _, m := range result {
		r := dto.MessagingResponse{
			Id:          uuid.UUID(m.Id),
			FromId:      uuid.UUID(m.FromId),
			RoomId:      uuid.UUID(m.RoomId),
			ContentType: m.ContentType,
			Contents:    m.Contents,
		}
		if m.RoomId != nil {
			r.RoomId = uuid.UUID(m.RoomId)
		}
		res = append(res, r)
	}
	return res, nil
}

func (s *Service) PublishMessaging(ctx context.Context, fromId uuid.UUID, toIdType string, toId uuid.UUID, contentType string, contents []string) (map[string]uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if contentType != "text" {
		exists, err1 := s.repository.HasFilepath(ctx, contentType+string(fromId[:]), contents)
		if err1 != nil {
			return nil, err1
		}
		if !exists {
			return nil, errors.New("bad request")
		}
		err = s.repository.RemoveFilepath(ctx, contentType+string(fromId[:]), contents)
		if err != nil {
			return nil, err
		}
	}
	p := payload.Marshal(payload.ChatMessage{
		Id:          id[:],
		FromId:      fromId[:],
		ToIdType:    toIdType,
		ToId:        toId[:],
		ContentType: contentType,
		Contents:    contents,
	})
	err = s.producer.PushMessage("chat.message", nil, p)
	if err != nil {
		slog.Error("fail to publish message", "err", err)
		return nil, err
	}
	return map[string]uuid.UUID{"id": id}, nil
}

func (s *Service) GeneratePresignedURL(ctx context.Context, id uuid.UUID, contentType string, n int) (res []dto.GeneratePresignedURLResponse, err error) {
	var filenames []string
	for range n {
		filename, err1 := uuid.NewV7()
		if err1 != nil {
			slog.Error("fail to generate uuid v7", "err", err1)
			return nil, err1
		}
		filenames = append(filenames, filename.String())
		res = append(res, dto.GeneratePresignedURLResponse{Filename: filename})
	}
	err = s.repository.SetFilepath(ctx, contentType+string(id[:]), filenames)
	if err != nil {
		slog.Error("fail to save filename", "err", err)
		return nil, err
	}
	for i := range n {
		p, err1 := s.presignClient.PresignPostObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(fmt.Sprintf("%s/%s", contentType, filenames[i])),
		}, func(opts *s3.PresignPostOptions) {
			opts.Expires = 1 * time.Hour
			opts.Conditions = []any{
				[]any{"content-length-range", 1, 1024 * 1024 * 1024},
				[]any{"starts-with", "$Content-Type", contentType},
			}
		})
		if err1 != nil {
			slog.Error("fail to generate presigned URL", "err", err)
			return nil, err
		}
		res[i].URL = p.URL
		res[i].Fields = p.Values
	}
	return res, nil
}

func (s *Service) GenerateSignedURL(ctx context.Context, memberId uuid.UUID, contentType string, filename uuid.UUID) (map[string]string, error) {
	ids, err := s.repository.FindIdsByFilename(ctx, gocql.UUID(filename))
	if err != nil {
		return nil, err
	}
	if !slices.Contains(ids, gocql.UUID(memberId)) {
		slog.Warn("this member id don't have authority to see file")
		return nil, errors.New("unauthorized request")
	}
	if contentType == "image" {
		contentType = "resized-image"
	}
	signedURL, err := s.signer.Sign(
		fmt.Sprintf("%s/%s/%s",
			s.cloudfrontURL, contentType, filename),
		time.Now().Add(1*time.Hour))
	if err != nil {
		slog.Error("fail to generate signed URL",
			"err", err)
		return nil, err
	}
	slog.Info("success to sign url", "signedURL", signedURL)
	return map[string]string{"url": signedURL}, nil
}
