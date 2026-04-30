package service

import (
	"backend/payload"
	pb "backend/proto"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Service) RelayMessage(ctx context.Context, id uuid.UUID, toIds [][]byte, roomId uuid.UUID, fromId uuid.UUID, contentType string, contents []string) error {
	var wg sync.WaitGroup
	relayToIdsByIP := make(map[string][][]byte)
	var rm sync.Mutex
	var pushToIds [][]byte
	var pm sync.Mutex
	for _, tid := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxr, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			ips, err := s.repository.GetServerIPs(ctxr, string(tid))
			if (err != nil || ips == nil) && !bytes.Equal(tid, fromId[:]) {
				pm.Lock()
				pushToIds = append(pushToIds, tid)
				pm.Unlock()
				return
			}
			rm.Lock()
			for _, ip := range ips {
				log.Printf("ip: %v", ip)
				relayToIdsByIP[ip] = append(relayToIdsByIP[ip], tid)
			}
			rm.Unlock()
		}()
	}
	wg.Wait()
	log.Printf("pushToId: %v, relayToId: %v", pushToIds, relayToIdsByIP)

	if pushToIds != nil {
		p, _ := json.Marshal(payload.PreparedMessage{
			ToIds:       pushToIds,
			RoomId:      roomId[:],
			FromId:      fromId[:],
			ContentType: contentType,
			Contents:    contents,
		})
		err := s.producer.PushMessage("notify_message", nil, p)
		if err != nil {
			return err
		}
		pushToIds = nil
	}

	for ip, tids := range relayToIdsByIP {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("relay start ip: %v, tids: %v", ip, tids)
			s.ccsMutex.RLock()
			cc, ok := s.clientConns[ip]
			s.ccsMutex.RUnlock()
			var err error
			if !ok {
				log.Printf("try to make connection...")
				var opts []grpc.DialOption
				opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
				cc, err = grpc.NewClient(ip+":50051", opts...)
				if err != nil {
					slog.Error("fail to get *ClientConn for relay",
						"err", err)
					pm.Lock()
					pushToIds = append(pushToIds, tids...)
					pm.Unlock()
					return
				}
				s.ccsMutex.Lock()
				s.clientConns[ip] = cc
				s.ccsMutex.Unlock()
			}
			client := pb.NewMessagingServiceClient(cc)
			req := pb.RelayMessagingRequest{
				Id:          id[:],
				ToIds:       tids,
				RoomId:      roomId[:],
				FromId:      fromId[:],
				ContentType: contentType,
				Contents:    contents,
			}
			ctxt, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			res, err := client.RelayMessaging(ctxt, &req)
			if err != nil {
				slog.Error("fail to relay messaging", "err", err)
				pm.Lock()
				pushToIds = append(pushToIds, tids...)
				pm.Unlock()
				//st, ok := status.FromError(err)
				//if ok && (st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded) {
				s.ccsMutex.Lock()
				err = s.clientConns[ip].Close()
				if err != nil {
					slog.Error("fail to close grpc client connection", "err", err)
				}
				delete(s.clientConns, ip)
				s.ccsMutex.Unlock()
				var wg1 sync.WaitGroup
				for _, tid := range tids {
					wg1.Add(1)
					go func() {
						defer wg1.Done()
						s.checkAndRemoveStaleIP(tid, ip)
					}()
				}
				wg1.Wait()
				//}
				return
			}

			if res == nil || len(res.PushToIds) == 0 {
				return
			}
			pm.Lock()
			pushToIds = append(pushToIds, res.PushToIds...)
			pm.Unlock()
			var wg1 sync.WaitGroup
			for _, tid := range res.PushToIds {
				wg1.Add(1)
				go func() {
					defer wg1.Done()
					s.checkAndRemoveStaleIP(tid, ip)
				}()
			}
			wg1.Wait()
		}()
	}
	wg.Wait()

	var filteredIds [][]byte
	for _, ptid := range pushToIds {
		if !bytes.Equal(ptid, fromId[:]) {
			filteredIds = append(filteredIds, ptid)
		}
	}

	if filteredIds != nil {
		p, _ := json.Marshal(payload.PreparedMessage{
			ToIds:       filteredIds,
			RoomId:      roomId[:],
			FromId:      fromId[:],
			ContentType: contentType,
			Contents:    contents,
		})
		err := s.producer.PushMessage("notify_message", nil, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) checkAndRemoveStaleIP(tid []byte, ip string) {
	ctxr, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()
	currentIPs, err1 := s.repository.GetServerIPs(ctxr, string(tid))
	if err1 != nil {
		return
	}
	if slices.Contains(currentIPs, ip) {
		err1 = s.repository.RemoveServerIP(ctxr, string(tid), ip)
		if err1 != nil {
			return
		}
	}
}
