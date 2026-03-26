package service

import (
	pb "backend/proto"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Service) PropagateSignal(ctx context.Context, fromId, toIdRaw string, signal json.RawMessage) error {
	toId, err := uuid.Parse(toIdRaw)
	if err != nil {
		slog.Error("fail to parse")
		return err
	}
	ip, err := s.repository.GetServerIP(ctx, toId)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := pb.RelaySignalRequest{
		FromId: fromId,
		ToId:   toIdRaw,
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
