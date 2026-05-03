package service

import (
	"backend/payload"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) PreprocessNotification(ctx context.Context, toIds [][]byte, roomId, fromId uuid.UUID, contentType string, content []string) error {
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

	p := payload.NotificationMessage{
		RoomName:   roomName,
		SenderName: senderName,
		Text:       content[0],
	}

	if contentType == "image" {
		imageURL, err1 := s.generateSignedURL(contentType, content[0])
		if err1 != nil {
			slog.Error("fail to generate Signed URL", "err", err1)
			return err1
		}
		p.ImageURL = imageURL
	}

	if contentType != "text" {
		p.Text = fmt.Sprintf("(%v)", contentType)
	}

	if len(fcmtm) > 0 {
		p.TokenMap = fcmtm
		err = s.producer.PushMessage("fcm_notification", payload.Marshal(p))
		if err != nil {
			return err
		}
	}
	if len(apntm) > 0 {
		p.TokenMap = apntm
		err = s.producer.PushMessage("apn_notification", payload.Marshal(p))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) generateSignedURL(contentType, filename string) (string, error) {
	if contentType == "image" {
		contentType = "resized-image"
	}
	signedURL, err := s.signer.Sign(
		fmt.Sprintf("%s/%s/%s",
			s.cloudfrontURL, contentType, filename),
		time.Now().Add(1*time.Hour))
	if err != nil {
		slog.Error("fail to generate signed URL",
			"err", err)
		return "", err
	}
	slog.Info("success to sign url", "signedURL", signedURL)
	return signedURL, nil
}
