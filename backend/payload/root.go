package payload

import (
	"encoding/json"
)

type OnlineConversationSignal struct {
	FromId json.RawMessage `json:"fromId"`
	ToId   json.RawMessage `json:"toId"`
	Signal json.RawMessage `json:"signal,omitempty"`
}

type ChatMessaging struct {
	ToIds       []json.RawMessage `json:"toIds"`
	RoomId      json.RawMessage   `json:"roomId,omitempty"`
	FromId      json.RawMessage   `json:"fromId"`
	ContentType string            `json:"contentType"`
	Content     string            `json:"content"`
}
