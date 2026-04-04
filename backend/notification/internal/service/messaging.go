package service

import (
	"backend/notification/internal/constant"
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

	var em sync.Mutex
	var es []error
	var wg sync.WaitGroup
	id, err0 := uuid.NewV7()
	if err0 != nil {
		slog.Error("fail to create uuid V7 for messaging", "err", err0)
		return err0
	}
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			err := s.repository.SaveMessaging(
				ctxt, gocql.UUID(id), gocql.UUID(toId), gocql.UUID(roomId), gocql.UUID(fromId),
				contentType, content)
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

	apntm := make(map[string]uuid.UUID)
	var am sync.Mutex
	fcmtm := make(map[string]uuid.UUID)
	var fm sync.Mutex
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			os, t, err := s.repository.FindDevicePushToken(ctxt, gocql.UUID(toId))
			if errors.Is(err, gocql.ErrNotFound) || t == "" {
				err = nil
				return
			}
			if err != nil {
				em.Lock()
				es = append(es, err)
				em.Unlock()
				return
			}
			if os == "ios" {
				am.Lock()
				apntm[t] = toId
				am.Unlock()
				return
			}
			fm.Lock()
			fcmtm[t] = toId
			fm.Unlock()
		}()
	}
	wg.Wait()
	if errors.Join(es...) != nil {
		return errors.Join(es...)
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
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	type Aps struct {
		Alert Alert `json:"alert"`
	}
	type Payload struct {
		Aps Aps `json:"aps"`
	}
	payload, _ := json.Marshal(Payload{
		Aps: Aps{
			Alert: Alert{
				Title: "",
				Body:  "",
			},
		},
	})
	var failedIds []uuid.UUID
	var mu sync.Mutex
	for t := range apntm {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err1 := http.NewRequestWithContext(ctx, "POST", constant.ApplePushURL+t, bytes.NewBuffer(payload))
			if err1 != nil {
				em.Lock()
				es = append(es, err1)
				em.Unlock()
				return
			}
			req.Header.Set("apns-topic", s.BundleIdentifier)
			req.Header.Set("authorization", authorization)
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			res, err1 := client.Do(req)
			if err1 != nil {
				slog.Error("fail to send request", "err1", err1)
				em.Lock()
				es = append(es, err1)
				em.Unlock()
				return
			}
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				return
			}
			mu.Lock()
			slog.Error("fail to send apn notification",
				"failedId", apntm[t])
			failedIds = append(failedIds, apntm[t])
			mu.Unlock()
		}()
	}
	wg.Wait()

	client, err := s.app.Messaging(ctx)
	if err != nil {
		slog.Error("fail to create fcm messaging client",
			"err", err)
		return err
	}

	var fcmts []string
	for t := range fcmtm {
		fcmts = append(fcmts, t)
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
	for _, failedId := range failedIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err2 := s.repository.RemoveInvalidDeviceToken(ctx, gocql.UUID(failedId))
			if err2 != nil {
				em.Lock()
				es = append(es, err2)
				em.Unlock()
			}
		}()
	}
	wg.Wait()
	if errors.Join(es...) != nil {
		return errors.Join(es...)
	}

	return nil
}
