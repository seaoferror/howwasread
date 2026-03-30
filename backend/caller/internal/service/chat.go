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
	toIds []uuid.UUID,
	roomId, fromId uuid.UUID,
	contentType, content string) error {
	var pushToIds [][]byte
	var relayToIdsByIPs map[string][][]byte
	for _, toId := range toIds {
		ctx := context.Background()
		ip, err := s.repository.GetServerIP(ctx, string(toId[:]))
		if err != nil {
			continue
		}
		if ip == "" {
			pushToIds = append(pushToIds, toId[:])
			continue
		}
		relayToIdsByIPs[ip] = append(relayToIdsByIPs[ip], toId[:])
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	for ip, ids := range relayToIdsByIPs {
		func() {
			cc, err := grpc.NewClient(ip+":50051", opts...)
			if err != nil {
				slog.Error("fail to get *ClientConn",
					"err", err)
				return
			}
			defer func(cc *grpc.ClientConn) {
				err = cc.Close()
				if err != nil {
					slog.Error("fail to close *ClientConn",
						"err", err)
				}
			}(cc)
			client := pb.NewMessagingServiceClient(cc)
			req := pb.RelayMessagingRequest{
				ToIds:       ids,
				FromId:      fromId[:],
				ContentType: contentType,
				Content:     content,
			}
			if roomId != uuid.Nil {
				req.RoomId = roomId[:]
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err = client.RelayMessaging(ctx, &req)
			if err != nil {
				slog.Error("fail to relay messaging",
					"err", err)
				return
			}
		}()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cc, err := grpc.NewClient("push-service:50051", opts...)
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
	client := pb.NewNotificationServiceClient(cc)
	req := pb.NotifyMessagingRequest{
		ToIds:       pushToIds,
		FromId:      fromId[:],
		ContentType: contentType,
		Content:     content,
	}
	if roomId != uuid.Nil {
		req.RoomId = roomId[:]
	}
	_, err = client.NotifyMessaging(ctx, &req)
	if err != nil {
		slog.Error("fail to relay messaging",
			"err", err)
		return err
	}
	return nil
}
