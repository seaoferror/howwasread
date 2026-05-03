package service

import (
	"backend/notification/internal/constant"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) NotifyMessaging(ctx context.Context, toIds [][]byte, roomId, fromId uuid.UUID, contentType string, content []string) error {
	log.Print("start notify message...")
	var em sync.Mutex
	var es []error
	var wg sync.WaitGroup
	apntm := make(map[string]uuid.UUIDs)
	var am sync.Mutex
	fcmtm := make(map[string]uuid.UUIDs)
	var fm sync.Mutex
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			result, err := s.repository.FindPushTokensById(ctxt, gocql.UUID(toId))
			if err != nil {
				em.Lock()
				es = append(es, err)
				em.Unlock()
				return
			}
			for _, d := range result {
				if d.OS == "ios" {
					am.Lock()
					apntm[d.DevicePushToken] = uuid.UUIDs{uuid.UUID(d.Id), uuid.UUID(d.DeviceId)}
					am.Unlock()
					return
				}
				fm.Lock()
				fcmtm[d.DevicePushToken] = uuid.UUIDs{uuid.UUID(d.Id), uuid.UUID(d.DeviceId)}
				fm.Unlock()
			}
		}()
	}
	wg.Wait()
	err0 := errors.Join(es...)
	if err0 != nil {
		return err0
	}

	body := content[0]
	if contentType != "text" {
		body = fmt.Sprintf("(%s)", contentType)
	}

	senderName, err := s.repository.FindNameById(ctx, gocql.UUID(fromId))
	if err != nil {
		return err
	}
	var roomName string
	if !bytes.Equal(toIds[0], roomId[:]) {
		roomName, err = s.repository.FindRoomNameById(ctx, gocql.UUID(roomId))
	}
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": s.teamId,
		"iat": time.Now().Unix(),
	})
	token.Header["kid"] = s.keyId

	authorization, err := token.SignedString(s.privateKey)
	if err != nil {
		slog.Error("fail to make authorization token", "err", err)
		return err
	}
	type Alert struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Body     string `json:"body"`
	}
	type Aps struct {
		Alert Alert `json:"alert"`
	}
	type Payload struct {
		Aps Aps `json:"aps"`
	}
	alert := Alert{
		Body: body,
	}
	alert.Title = senderName
	if roomName != "" {
		alert.Title = roomName
		alert.Subtitle = senderName
	}

	payload, _ := json.Marshal(Payload{
		Aps: Aps{
			Alert: alert,
		},
	})
	var failedIds []uuid.UUIDs
	var mu sync.Mutex
	for t := range apntm {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err1 := http.NewRequestWithContext(ctx, http.MethodPost, constant.ApplePushURL+t, bytes.NewBuffer(payload))
			if err1 != nil {
				em.Lock()
				es = append(es, err1)
				em.Unlock()
				return
			}
			req.Header.Set("apns-topic", s.BundleIdentifier)
			req.Header.Set("authorization", "bearer "+authorization)
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			res, err1 := client.Do(req)
			defer res.Body.Close()
			bodyRaw, err2 := io.ReadAll(res.Body)
			if err2 != nil {
				slog.Error("fail to read body", "err", err2)
			}
			log.Printf("apn status code: %v", res.StatusCode)
			if res.StatusCode == http.StatusOK {
				slog.Info("success to notify message to apn",
					"body", string(bodyRaw))
				return
			}
			if err1 != nil {
				slog.Error("fail to send request", "err1", err1)
				em.Lock()
				es = append(es, err1)
				em.Unlock()
			}
			log.Printf("failed apn body: %v", string(bodyRaw))
			mu.Lock()
			slog.Error("fail to send apn notification",
				"failedId", apntm[t])
			failedIds = append(failedIds, apntm[t])
			mu.Unlock()
		}()
	}
	wg.Wait()

	var fcmts []string
	for t := range fcmtm {
		fcmts = append(fcmts, t)
	}

	if len(fcmts) > 0 {
		client, err := s.app.Messaging(ctx)
		if err != nil {
			slog.Error("fail to create fcm messaging client",
				"err", err)
			return err
		}

		message := &messaging.MulticastMessage{
			Notification: &messaging.Notification{
				Title: senderName,
				Body:  body,
			},
			Tokens: fcmts,
		}
		br, err := client.SendEachForMulticast(ctx, message)
		if err != nil {
			return err
		}
		if br.FailureCount > 0 {
			for i, resp := range br.Responses {
				if !resp.Success {
					slog.Info("fail to send fcm notification",
						"failedId", fcmtm[fcmts[i]])
					failedIds = append(failedIds, fcmtm[fcmts[i]])
				}
			}
		}
	}
	for _, failedId := range failedIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err2 := s.repository.RemovePushTokenByIdAndDeviceId(
				ctx, gocql.UUID(failedId[0]), gocql.UUID(failedId[1]))
			if err2 != nil {
				em.Lock()
				es = append(es, err2)
				em.Unlock()
			}
		}()
	}
	wg.Wait()
	err0 = errors.Join(es...)
	if err0 != nil {
		return err0
	}

	return nil
}
