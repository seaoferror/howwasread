package service

import (
	pb "backend/proto"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (s *Service) PropagateSignal(ctx context.Context, toIds [][]byte, fromId []byte, signal json.RawMessage) {
	var wg sync.WaitGroup
	relayToIdsByIPs := make(map[string][][]byte)
	var rm sync.Mutex
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ip, err := s.repository.GetServerIP(ctx, string(toId))
			if err != nil || ip == "" {
				return
			}
			rm.Lock()
			relayToIdsByIPs[ip] = append(relayToIdsByIPs[ip], toId)
			rm.Unlock()
		}()
	}
	wg.Wait()

	for ip, ids := range relayToIdsByIPs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			s.ccsMutex.RLock()
			cc, ok := s.clientConns[ip]
			s.ccsMutex.RUnlock()
			if !ok {
				var opts []grpc.DialOption
				opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
				cc, err = grpc.NewClient(ip+":50051", opts...)
				if err != nil {
					slog.Error("fail to make new grpc client", "err", err)
					return
				}
				s.ccsMutex.Lock()
				s.clientConns[ip] = cc
				s.ccsMutex.Unlock()
			}
			client := pb.NewSignalServiceClient(cc)
			req := pb.RelaySignalRequest{
				ToIds:  ids,
				FromId: fromId,
				Signal: signal,
			}
			ctxg, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			_, err = client.RelaySignal(ctxg, &req)
			st, ok := status.FromError(err)
			if ok && (st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded) {
				s.ccsMutex.RLock()
				err = s.clientConns[ip].Close()
				s.ccsMutex.RUnlock()
				if err != nil {
					slog.Error("fail to close grpc client connection", "err", err)
				}
				s.ccsMutex.Lock()
				delete(s.clientConns, ip)
				s.ccsMutex.Unlock()
				var wg1 sync.WaitGroup
				for _, tid := range ids {
					wg1.Add(1)
					go func() {
						defer wg1.Done()
						currentIP, err1 := s.repository.GetServerIP(ctx, string(tid))
						if err1 != nil {
							return
						}
						if currentIP == ip {
							err1 = s.repository.RemoveServerIP(ctx, string(tid))
							if err1 != nil {
								return
							}
						}
					}()
				}
				wg1.Wait()
			}
			if err != nil {
				slog.Error("fail to relay signal", "err", err)
			}
		}()
	}
	wg.Wait()

	return
}
