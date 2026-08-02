package service

import (
	"backend/common/payload"
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

func (s *Service) PreprocessNotification(
	ctx context.Context,
	notificationId uint8,
	messageId uuid.UUID,
	toIds [][]byte,
	roomId, fromId uuid.UUID,
	contentType string,
	content []string,
) {
	log.Print("start notify message...")
	var em sync.Mutex
	var es []error
	var wg sync.WaitGroup
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
					apntm[d.DevicePushToken] = uuid.UUID(d.Id)
					am.Unlock()
					return
				}
				fm.Lock()
				fcmtm[d.DevicePushToken] = uuid.UUID(d.Id)
				fm.Unlock()
			}
		}()
	}
	wg.Wait()
	err0 := errors.Join(es...)
	if err0 != nil {
		return
	}

	senderName, err := s.repository.FindNameById(ctx, gocql.UUID(fromId))
	if err != nil {
		return
	}
	var roomName string
	if !bytes.Equal(toIds[0], roomId[:]) {
		roomName, err = s.repository.FindRoomNameById(ctx, gocql.UUID(roomId))
	}
	if err != nil {
		return
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
			return
		}
		p.ImageURL = imageURL
	}

	if contentType != "text" {
		p.Text = fmt.Sprintf("(%v)", contentType)
	}

	kafkaKey := append(messageId[:], notificationId)
	if len(fcmtm) > 0 {
		p.TokenMap = fcmtm
		s.producer.PushMessage("fcm-notification", kafkaKey, payload.Marshal(p), nil)
	}
	if len(apntm) > 0 {
		p.TokenMap = apntm
		s.producer.PushMessage("apn-notification", kafkaKey, payload.Marshal(p), nil)
	}
	return
}

func (s *Service) generateSignedURL(contentType, filename string) (string, error) {
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
