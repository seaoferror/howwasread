package payload

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func Marshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		slog.Error("fail to marshal", "err", err)
		return nil
	}
	return b
}

type OnlineConversationSignal struct {
	FromId []byte          `json:"fromId"`
	ToIds  [][]byte        `json:"toIds"`
	Signal json.RawMessage `json:"signal,omitempty"`
}

type ChatMessage struct {
	Id          []byte   `json:"id"`
	FromId      []byte   `json:"fromId"`
	ToIdType    string   `json:"toIdType"`
	ToId        []byte   `json:"toId"`
	ContentType string   `json:"contentType"`
	Contents    []string `json:"contents"`
}

type PreparedMessage struct {
	NotificationId uint8    `json:"notificationId,omitempty"`
	Id             []byte   `json:"id"`
	ToIds          [][]byte `json:"toIds"`
	RoomId         []byte   `json:"roomId"`
	FromId         []byte   `json:"fromId"`
	ContentType    string   `json:"contentType"`
	Contents       []string `json:"contents"`
}

type NotificationMessage struct {
	TokenMap   map[string]uuid.UUID `json:"tokenMap"`
	RoomName   string               `json:"roomName,omitempty"`
	SenderName string               `json:"senderName"`
	Text       string               `json:"text"`
	ImageURL   string               `json:"imageURL,omitempty"`
}

type ConversationRequest struct {
	Id uuid.UUID `json:"id"`
}

type OnlineConversationNotification struct {
	ConversationId uuid.UUID `json:"conversationId"`
	MemberId       uuid.UUID `json:"memberId"`
	ScheduledTime  time.Time `json:"scheduledTime"`
}
