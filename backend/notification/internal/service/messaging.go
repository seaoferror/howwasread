package service

import (
	"backend/notification/internal/constant"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) NotifyMessaging(ctx context.Context,
	toIds []uuid.UUID, roomId, fromId uuid.UUID,
	contentType, content string) error {

	var wg sync.WaitGroup
	ec := make(chan error)

	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			id, err := uuid.NewV7()
			if err != nil {
				slog.Error("fail to create uuid V7 for messaging", "err", err)
				ec <- err
				return
			}
			err = s.repository.SaveMessaging(
				ctxt, gocql.UUID(id), gocql.UUID(toId), gocql.UUID(roomId), gocql.UUID(fromId),
				contentType, content)
			if err != nil {
				ec <- err
			}
		}()
	}
	wg.Wait()

	apntc := make(chan map[uuid.UUID]string)
	fcmtc := make(chan map[uuid.UUID]string)
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			os, t, err := s.repository.FindDevicePushToken(ctxt, gocql.UUID(toId))
			if err != nil {
				ec <- err
				return
			}
			if os == "ios" {
				apntc <- map[uuid.UUID]string{toId: t}
				return
			}
			fcmtc <- map[uuid.UUID]string{toId: t}
		}()
	}
	wg.Wait()
	close(apntc)
	close(fcmtc)
	apntm := make(map[string]uuid.UUID)
	var apnts []string
	for t := range fcmtc {
		for k, v := range t {
			apntm[v] = k
			apnts = append(apnts, v)
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": s.teamId,
		"iat": time.Now().Unix(),
	})
	token.Header["kid"] = s.keyId
	func() {
		authorization, err := token.SignedString(s.privateKey)
		if err != nil {
			slog.Error("fail to make authorization token", "err", err)
			ec <- err
			return
		}
		type Alert struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		}
		type Aps struct {
			Alert Alert `json:"alert"`
		}
		type Payload struct {
			Aps Aps `json:"aps"`
		}
		payload, err := json.Marshal(Payload{
			Aps: Aps{
				Alert: Alert{
					Title: "",
					Body:  "",
				},
			},
		})
		if err != nil {
			slog.Error("fail to marshal payload for apn",
				"err", err)
			ec <- err
			return
		}
		for _, t := range apnts {
			go func() {
				req, err1 := http.NewRequest("POST", constant.ApplePushURL+t, bytes.NewBuffer(payload))
				if err1 != nil {
					ec <- err1
					return
				}
				req.Header.Set("apns-topic", s.BundleIdentifier)
				req.Header.Set("authorization", authorization)
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				res, err1 := client.Do(req)
				if err1 != nil {
					slog.Error("fail to send request", "err1", err1)
					ec <- err1
					return
				}
				defer res.Body.Close()
				bodyRaw, err1 := io.ReadAll(res.Body)
				if err1 != nil {
					slog.Error("fail to read body", "err", err1)
					ec <- err1
					return
				}
				if res.StatusCode == http.StatusOK {
					return
				}
				slog.Error("fail to send apn notification",
					"failedId", apntm[t],
					"statusCode", res.StatusCode,
					"response", string(bodyRaw))
			}()
		}
	}()

	fcmtm := make(map[string]uuid.UUID)
	var fcmts []string
	for t := range fcmtc {
		for k, v := range t {
			fcmtm[v] = k
			fcmts = append(fcmts, v)
		}
	}

	func() {
		client, err := s.app.Messaging(ctx)
		if err != nil {
			slog.Error("fail to create fcm client",
				"err", err)
			ec <- err
			return
		}

		message := &messaging.MulticastMessage{
			Notification: &messaging.Notification{
				Title: "",
				Body:  "",
			},
			Tokens: fcmts,
		}
		br, err := client.SendEachForMulticast(ctx, message)
		if err != nil {
			ec <- err
			return
		}
		if br.FailureCount > 0 {
			var failedIds []string
			for i, resp := range br.Responses {
				if !resp.Success {
					// The order of responses corresponds to the order of the registration tokens.
					failedIds = append(failedIds, fcmtm[fcmts[i]].String())
				}
			}
			slog.Error("failed ids to get push notification", "failedIds", failedIds)
			//TODO: publish this to kafka and retry?
		}
	}()

	close(ec)
	var errs []error
	for err := range ec {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
