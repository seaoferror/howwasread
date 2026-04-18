package service

import (
	pb "backend/proto"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"slices"
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
	roomId := toId
	if toIdType == "personal" {
		toIds = append(toIds, toId[:])
		toIds = append(toIds, fromId[:])
	}
	if toIdType == "room" {
		participantIds, err := s.repository.FindParticipantIds(ctx, gocql.UUID(toId))
		if err != nil {
			return err
		}
		for _, pid := range participantIds {
			toIds = append(toIds, pid[:])
		}
	}
	if contentType != "text" {
		filename, err := uuid.Parse(content)
		if err != nil {
			slog.Error("fail to parse uuid from content", "err", err)
			return err
		}
		var refinedToIds []gocql.UUID
		for _, tid := range toIds {
			refinedToIds = append(refinedToIds, gocql.UUID(tid))
		}
		err = s.repository.SaveIdsByFileName(ctx, refinedToIds, gocql.UUID(filename))
		if err != nil {
			return err
		}
	}
	var wg sync.WaitGroup
	var es []error
	var em sync.Mutex
	for _, tid := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if bytes.Equal(tid, roomId[:]) {
				err := s.repository.SaveMessaging(ctx, gocql.UUID(id), gocql.UUID(tid), gocql.UUID(fromId), gocql.UUID(fromId), contentType, content)
				if err != nil {
					em.Lock()
					es = append(es, err)
					em.Unlock()
				}
				return
			}
			err := s.repository.SaveMessaging(ctx, gocql.UUID(id), gocql.UUID(tid), gocql.UUID(fromId), gocql.UUID(roomId), contentType, content)
			if err != nil {
				em.Lock()
				es = append(es, err)
				em.Unlock()
			}
		}()
	}
	wg.Wait()
	err0 := errors.Join(es...)
	if err0 != nil {
		return err0
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
				ctxr, cancel := context.WithTimeout(ctx, 1*time.Second)
				defer cancel()
				defer wg.Done()
				ips, err := s.repository.GetServerIPs(ctxr, string(tid))
				if err != nil || ips == nil {
					pm.Lock()
					pushToIds = append(pushToIds, tid)
					pm.Unlock()
					return
				}
				rm.Lock()
				for _, ip := range ips {
					relayToIdsByIP[ip] = append(relayToIdsByIP[ip], tid)
				}
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
				var err error
				if !ok {
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
					Content:     content,
				}
				ctxt, cancel := context.WithTimeout(ctx, time.Second*5)
				defer cancel()
				res, err := client.RelayMessaging(ctxt, &req)
				st, ok := status.FromError(err)
				if ok && (st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded) {
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
				pushToIds = append(pushToIds, res.PushToIds...)
				var wg1 sync.WaitGroup
				for _, tid := range res.PushToIds {
					wg1.Add(1)
					go func() {
						defer wg1.Done()
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

		s.ccsMutex.RLock()
		cc := s.clientConns["push-service"]
		s.ccsMutex.RUnlock()

		client := pb.NewNotificationServiceClient(cc)
		req := pb.NotifyMessagingRequest{
			ToIds:       filteredIds,
			FromId:      fromId[:],
			ContentType: contentType,
			Content:     content,
		}
		if toIdType == "room" {
			req.RoomId = roomId[:]
		}
		ctxt, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_, err := client.NotifyMessaging(ctxt, &req)
		if err != nil {
			slog.Error("fail to notify message", "err", err)
		}
		return
	}()

	return nil
}
