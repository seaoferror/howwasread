package payload

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
)

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
	Id          []byte   `json:"id,omitempty"`
	ToIds       [][]byte `json:"toIds"`
	RoomId      []byte   `json:"roomId"`
	FromId      []byte   `json:"fromId"`
	ContentType string   `json:"contentType"`
	Contents    []string `json:"contents"`
}

type NotificationMessage struct {
	TokenMap   map[string]uuid.UUIDs `json:"tokenMap"`
	RoomName   string                `json:"roomName,omitempty"`
	SenderName string                `json:"senderName"`
	Text       string                `json:"text"`
	ImageURL   string                `json:"imageURL,omitempty"`
}

func Marshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		slog.Error("fail to marshal", "err", err)
		return nil
	}
	return b
}
