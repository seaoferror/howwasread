package service

import (
	pb "backend/proto"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Service) PropagateSignal(fromId, toId uuid.UUID, signal json.RawMessage) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ip, err := s.repository.GetServerIP(ctx, string(toId[:]))
	if err != nil {
		return err
	}
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

	client := pb.NewSignalServiceClient(cc)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := pb.RelaySignalRequest{
		FromId: fromId[:],
		ToId:   toId[:],
		Signal: signal,
	}
	_, err = client.RelaySignal(ctx, &req)
	if err != nil {
		slog.Error("fail to relay signal",
			"err", err)
		return err
	}
	return nil
}
