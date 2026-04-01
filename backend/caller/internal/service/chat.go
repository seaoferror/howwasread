package service

import (
	pb "backend/proto"
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (s *Service) PropagateMessaging(ctx context.Context, toIds []uuid.UUID, roomId []byte, fromId uuid.UUID, contentType, content string) error {
	ec := make(chan error)
	var pushToIds [][]byte
	relayToIdsByIPs := make(map[string][][]byte)
	for _, toId := range toIds {
		func() {
			ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			ip, err := s.repository.GetServerIP(ctxt, string(toId[:]))
			if err != nil {
				ec <- err
				return
			}
			if ip == "" {
				pushToIds = append(pushToIds, toId[:])
				return
			}
			relayToIdsByIPs[ip] = append(relayToIdsByIPs[ip], toId[:])
		}()
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for ip, ids := range relayToIdsByIPs {
		wg.Add(1)
		go func() {
			var err error
			defer wg.Done()
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
					ec <- err
					return
				}
			}
			cc := s.clientConns[ip]
			s.ccsMutex.RUnlock()

			client := pb.NewMessagingServiceClient(cc)
			req := pb.RelayMessagingRequest{
				ToIds:       ids,
				RoomId:      roomId,
				FromId:      fromId[:],
				ContentType: contentType,
				Content:     content,
			}
			ctxt, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			_, err = client.RelayMessaging(ctxt, &req)
			if err != nil {
				slog.Error("fail to relay messaging", "err", err)
				st, ok := status.FromError(err)
				if ok && (st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded) {
					s.ccsMutex.Lock()
					err = s.clientConns[ip].Close()
					if err != nil {
						slog.Error("fail to close grpc client connection", "err", err)
						ec <- err
					}
					delete(s.clientConns, ip)
					s.ccsMutex.Unlock()
				}
				mu.Lock()
				pushToIds = append(pushToIds, ids...)
				mu.Unlock()
				ec <- err
				return
			}
		}()
	}
	wg.Wait()

	func() {
		ctxt, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		var err error
		s.ccsMutex.RLock()
		if s.clientConns["push-service"] == nil {
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
			s.ccsMutex.Lock()
			s.clientConns["push-service"], err = grpc.NewClient("push-service:50051", opts...)
			s.ccsMutex.Unlock()
			if err != nil {
				slog.Error("fail to get *ClientConn",
					"err", err)
				ec <- err
				return
			}
		}
		cc := s.clientConns["push-service"]
		s.ccsMutex.RUnlock()

		client := pb.NewNotificationServiceClient(cc)
		req := pb.NotifyMessagingRequest{
			ToIds:       pushToIds,
			RoomId:      roomId,
			FromId:      fromId[:],
			ContentType: contentType,
			Content:     content,
		}
		_, err = client.NotifyMessaging(ctxt, &req)
		if err != nil {
			slog.Error("fail to relay messaging",
				"err", err)
			ec <- err
			return
		}
	}()

	close(ec)
	var errs []error
	for err := range ec {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
