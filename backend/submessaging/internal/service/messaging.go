package service

import (
	pb "backend/proto"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// ManageMessaging return error which will stop kafka message consumption
// in case which fail to persistent message in C*
// redis, fcm, apn, grpc error will not stop consumption
func (s *Service) ManageMessaging(ctx context.Context, id uuid.UUID, fromId uuid.UUID, toIdType string, toId uuid.UUID, contentType string, content string) error {
	var toIds [][]byte
	var roomId []byte
	if toIdType == "personal" {
		toIds = append(toIds, toId[:])
		toIds = append(toIds, fromId[:])
	}
	if toIdType == "room" {
		roomId = toId[:]
		participantIds, err := s.repository.FindParticipantIds(ctx, gocql.UUID(toId))
		if err != nil {
			return err
		}
		for _, pid := range participantIds {
			toIds = append(toIds, pid[:])
		}
	}
	var wg sync.WaitGroup
	var es []error
	var em sync.Mutex
	for _, tid := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.repository.SaveMessaging(ctx, gocql.UUID(id), gocql.UUID(tid), gocql.UUID(fromId), roomId, contentType, content)
			if err != nil {
				em.Lock()
				es = append(es, err)
				em.Unlock()
			}
		}()
	}
	wg.Wait()
	if errors.Join(es...) != nil {
		return errors.Join(es...)
	}

	//this will be background job
	go func() {
		ctx := context.Background()
		relayToIdsByIP := make(map[string][][]byte)
		var rm sync.Mutex
		var pushToIds [][]byte
		var pm sync.Mutex
		for _, tid := range toIds {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ip, err := s.repository.GetServerIP(ctx, string(tid))
				if err != nil {
					pm.Lock()
					pushToIds = append(pushToIds, tid)
					pm.Unlock()
					return
				}
				if ip == "" {
					pm.Lock()
					pushToIds = append(pushToIds, tid)
					pm.Unlock()
					return
				}
				rm.Lock()
				relayToIdsByIP[ip] = append(relayToIdsByIP[ip], tid)
				rm.Unlock()
			}()
		}
		wg.Wait()

		for ip, tids := range relayToIdsByIP {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.ccsMutex.RLock()
				cc, ok := s.clientConns[ip]
				s.ccsMutex.RUnlock()
				if !ok {
					var opts []grpc.DialOption
					opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
					cc, err := grpc.NewClient(ip+":50051", opts...)
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
					RoomId:      roomId,
					FromId:      fromId[:],
					ContentType: contentType,
					Content:     content,
				}
				ctxt, cancel := context.WithTimeout(ctx, time.Second*5)
				defer cancel()
				_, err := client.RelayMessaging(ctxt, &req)
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
					for _, tid := range tids {
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
					slog.Error("fail to relay messaging", "err", err)
					pm.Lock()
					pushToIds = append(pushToIds, tids...)
					pm.Unlock()
					return
				}
			}()
		}
		wg.Wait()

		var filteredIds [][]byte
		for _, ptid := range pushToIds {
			if !bytes.Equal(ptid, fromId[:]) {
				filteredIds = append(filteredIds, ptid)
			}
		}

		ctxt, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		s.ccsMutex.RLock()
		cc := s.clientConns["push-service"]
		s.ccsMutex.RUnlock()

		client := pb.NewNotificationServiceClient(cc)
		req := pb.NotifyMessagingRequest{
			ToIds:       filteredIds,
			RoomId:      roomId,
			FromId:      fromId[:],
			ContentType: contentType,
			Content:     content,
		}
		_, err := client.NotifyMessaging(ctxt, &req)
		if err != nil {
			slog.Error("fail to notify message", "err", err)
		}
		return
	}()

	return nil
}
