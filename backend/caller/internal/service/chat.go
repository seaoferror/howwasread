package service

import (
	pb "backend/proto"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Service) PropagateMessaging(
	ctx context.Context,
	toIds []uuid.UUID,
	roomId, fromId uuid.UUID,
	contentType, content string) error {
	for _, toId := range toIds {
		ip, err := s.repository.GetServerIP(ctx, string(toId[:]))
		if err != nil {
			return err
		}
		if ip != "" {
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
			cc, err := grpc.NewClient(ip+":50051", opts...)
			if err != nil {
				slog.Error("fail to get *ClientConn",
					"err", err)
				return err
			}
			defer func(cc *grpc.ClientConn) {
				err = cc.Close()
				if err != nil {
					slog.Error("fail to close *ClientConn",
						"err", err)
				}
			}(cc)

			client := pb.NewMessagingServiceClient(cc)
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			req := pb.RelayMessagingRequest{
				ToId:        toId[:],
				FromId:      fromId[:],
				ContentType: contentType,
				Content:     content,
			}
			if roomId != uuid.Nil {
				req.RoomId = roomId[:]
			}
			_, err = client.RelayMessaging(ctx, &req)
			if err != nil {
				slog.Error("fail to relay messaging",
					"err", err)
				return err
			}
			return nil
		}
		//TODO: send push

	}
	return nil
}
