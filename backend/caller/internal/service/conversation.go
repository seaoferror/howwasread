package service

import (
	pb "backend/proto"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (s *Service) PropagateSignal(ctx context.Context, fromId, toId uuid.UUID, signal json.RawMessage) error {
	ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	ip, err := s.repository.GetServerIP(ctxt, string(toId[:]))
	if err != nil {
		return err
	}
	s.ccsMutex.RLock()
	if s.clientConns[ip] == nil {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		s.ccsMutex.Lock()
		s.clientConns[ip], err = grpc.NewClient(ip+":50051", opts...)
		s.ccsMutex.Unlock()
		if err != nil {
			slog.Error("fail to get *ClientConn",
				"err", err)
			return err
		}
	}
	cc := s.clientConns[ip]
	s.ccsMutex.RUnlock()

	client := pb.NewSignalServiceClient(cc)

	req := pb.RelaySignalRequest{
		FromId: fromId[:],
		ToId:   toId[:],
		Signal: signal,
	}
	ctxt, cancel = context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = client.RelaySignal(ctxt, &req)
	if err != nil {
		slog.Error("fail to relay signal", "err", err)
		st, ok := status.FromError(err)
		if ok && (st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded) {
			s.ccsMutex.Lock()
			err = s.clientConns[ip].Close()
			delete(s.clientConns, ip)
			s.ccsMutex.Unlock()
			if err != nil {
				slog.Error("fail to close grpc client connection", "err", err)
				return err
			}
		}
		return err
	}
	return nil
}
