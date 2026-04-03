package service

import (
	"backend/payload"
	pb "backend/proto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (s *Service) PropagateMessaging(ctx context.Context, toIds [][]byte, roomId, fromId []byte, contentType, content string) error {
	var pushToIds [][]byte
	var pm sync.Mutex
	var wg sync.WaitGroup
	relayToIdsByIPs := make(map[string][][]byte)
	var rm sync.Mutex
	var errs []error
	var em sync.Mutex
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			ip, err := s.repository.GetServerIP(ctxt, string(toId))
			if err != nil || ip == "" {
				pm.Lock()
				pushToIds = append(pushToIds, toId)
				pm.Unlock()
				em.Lock()
				errs = append(errs, err)
				em.Unlock()
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
					slog.Error("fail to get *ClientConn",
						"err", err)
					pm.Lock()
					pushToIds = append(pushToIds, ids...)
					pm.Unlock()
					em.Lock()
					errs = append(errs, err)
					em.Unlock()
					return
				}
				s.ccsMutex.Lock()
				s.clientConns[ip] = cc
				s.ccsMutex.Unlock()
			}
			client := pb.NewMessagingServiceClient(cc)
			req := pb.RelayMessagingRequest{
				ToIds:       ids,
				RoomId:      roomId,
				FromId:      fromId,
				ContentType: contentType,
				Content:     content,
			}
			ctxt, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			out, err := client.RelayMessaging(ctxt, &req)
			if err != nil {
				slog.Error("fail to relay whole messaging", "err", err)
				em.Lock()
				errs = append(errs, err)
				em.Unlock()
				pm.Lock()
				pushToIds = append(pushToIds, out.FailedIds...)
				pm.Unlock()
			}
			st, ok := status.FromError(err)
			if ok && (st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded) {
				s.ccsMutex.RLock()
				err = s.clientConns[ip].Close()
				s.ccsMutex.RUnlock()
				if err != nil {
					slog.Error("fail to close grpc client connection", "err", err)
					errs = append(errs, err)
				}
				s.ccsMutex.Lock()
				delete(s.clientConns, ip)
				s.ccsMutex.Unlock()
			}
		}()
	}
	wg.Wait()

	ctxt, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	s.ccsMutex.RLock()
	cc := s.clientConns["push-service"]
	s.ccsMutex.RUnlock()

	client := pb.NewNotificationServiceClient(cc)
	req := pb.NotifyMessagingRequest{
		ToIds:       pushToIds,
		RoomId:      roomId,
		FromId:      fromId,
		ContentType: contentType,
		Content:     content,
	}
	_, err := client.NotifyMessaging(ctxt, &req)
	if err != nil {
		slog.Error("fail to notify messaging",
			"err", err)
		errs = append(errs, err)
		p, _ := json.Marshal(payload.ChatMessaging{
			ToIds:       pushToIds,
			RoomId:      roomId,
			FromId:      fromId,
			ContentType: contentType,
			Content:     content,
		})
		err = s.producer.PushDeadLetter(errors.Join(errs...), "chat.messaging", p)
		if err != nil {
			errs = append(errs, err)
		}
		return errors.Join(errs...)
	}
	return nil
}
