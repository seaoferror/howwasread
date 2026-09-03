package service

import (
	"backend/common/payload"
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

func (s *Service) PreprocessMessageNotification(
	ctx context.Context,
	notificationId uint8,
	messageId uuid.UUID,
	toIds [][]byte,
	roomId, fromId uuid.UUID,
	contentType string,
	content []string,
) {
	log.Print("start notify message...")
	tIds := make([]uuid.UUID, len(toIds))
	for _, toId := range toIds {
		tIds = append(tIds, uuid.UUID(toId))
	}
	apntm, fcmtm, err := s.getEachTokenMap(ctx, tIds)
	if err != nil {
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
		Title: senderName,
		Text:  content[0],
	}
	if roomName != "" {
		p.Title = roomName
		p.SubTitle = senderName
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
