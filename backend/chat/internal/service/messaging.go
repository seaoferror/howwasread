package service

import (
	"backend/chat/internal/dto"
	"backend/payload"
	"context"
	"encoding/json"
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
	result, err := s.repository.FindMessagesByToIdAndId(ctx, gocql.UUID(id), gocql.UUID(cursor))
	if err != nil {
		return nil, err
	}
	for _, m := range result {
		r := dto.MessagingResponse{
			Id:          uuid.UUID(m.Id),
			FromId:      uuid.UUID(m.FromId),
			RoomId:      uuid.UUID(m.RoomId),
			ContentType: m.ContentType,
			Content:     m.Content,
		}
		if m.RoomId != nil {
			r.RoomId = uuid.UUID(m.RoomId)
		}
		res = append(res, r)
	}
	return res, nil
}

func (s *Service) PublishMessaging(ctx context.Context, fromId uuid.UUID, toIdType string, toId uuid.UUID, contentType, content string) (map[string]uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if contentType != "text" {
		exists, err1 := s.repository.HasFilepath(ctx, string(fromId[:]), contentType+content)
		if err1 != nil {
			return nil, err1
		}
		if !exists {
			return nil, errors.New("bad request")
		}
		err1 = s.repository.RemoveFilepath(ctx, string(fromId[:]), contentType+content)
		if err1 != nil {
			return nil, err1
		}
	}
	p, _ := json.Marshal(payload.ChatMessage{
		Id:          id[:],
		FromId:      fromId[:],
		ToIdType:    toIdType,
		ToId:        toId[:],
		ContentType: contentType,
		Content:     content,
	})
	err = s.producer.PushMessage("chat.message", p)
	if err != nil {
		slog.Error("fail to publish message", "err", err)
		return nil, err
	}
	return map[string]uuid.UUID{"id": id}, nil
}

func (s *Service) GeneratePresignedURL(ctx context.Context, id uuid.UUID, contentType string) (*dto.GeneratePresignedURLResponse, error) {
	filename, err := uuid.NewV7()
	if err != nil {
		slog.Error("fail to generate uuid v7", "err", err)
		return nil, err
	}
	err = s.repository.SetFilepath(ctx, string(id[:]), contentType+filename.String())
	if err != nil {
		slog.Error("fail to save filename", "err", err)
		return nil, err
	}
	p, err := s.presignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("chat"),
		Key:    aws.String(fmt.Sprintf("%s/%s", contentType, filename.String())),
	}, func(opts *s3.PresignPostOptions) {
		opts.Expires = 1 * time.Second
		opts.Conditions = []any{
			[]any{"content-length-range", 1, 1024 * 1024 * 1024},
			[]any{"starts-with", "$Content-Type", contentType},
		}
	})
	if err != nil {
		slog.Error("fail to generate presigned URL", "err", err)
		return nil, err
	}
	res := dto.GeneratePresignedURLResponse{
		Filename: filename,
		URL:      p.URL,
		Fields:   p.Values,
	}
	return &res, nil
}

func (s *Service) GenerateSignedURL(ctx context.Context, id uuid.UUID, contentType string, filename uuid.UUID) (map[string]string, error) {
	ids, err := s.repository.FindIdsByFilename(ctx, filename)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(ids, gocql.UUID(id)) {
		slog.Warn("this member id don't have authority to see file")
		return nil, errors.New("unauthorized request")
	}
	signedURL, err := s.signer.Sign(
		fmt.Sprintf("https://d111111abcdef8.cloudfront.net/%s/%s",
			contentType, filename),
		time.Now().Add(1*time.Hour))
	if err != nil {
		slog.Error("fail to generate signed URL",
			"err", err)
		return nil, err
	}
	return map[string]string{"url": signedURL}, nil
}
