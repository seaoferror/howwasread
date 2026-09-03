package service

import (
	"backend/common/payload"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) PreprocessScheduledNotification(ctx context.Context, partitionId uuid.UUID, notifications map[uuid.UUID]map[int]string, contents map[int]string) {
	memberIds := make([]uuid.UUID, 0, len(notifications))
	for memberId, _ := range notifications {
		memberIds = append(memberIds, memberId)
	}
	apntm, fcmtm, err := s.getEachTokenMap(ctx, memberIds)
	if err != nil {
		return
	}
	p := payload.NotificationMessage{
		TokenMap: fcmtm,
		Title:    fmt.Sprintf("The online conversation about %s starts soon", contents[0]),
		Text:     "You can now enter the conversation room and talk!",
	}
	if len(fcmtm) > 0 {
		p.TokenMap = fcmtm
		s.producer.PushMessage("fcm-notification", partitionId[:], payload.Marshal(p), nil)
	}
	if len(apntm) > 0 {
		p.TokenMap = apntm
		s.producer.PushMessage("apn-notification", partitionId[:], payload.Marshal(p), nil)
	}
}
