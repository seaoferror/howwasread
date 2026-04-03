package payload

import "encoding/json"

type OnlineConversationSignal struct {
	FromId []byte          `json:"fromId"`
	ToIds  [][]byte        `json:"toIds"`
	Signal json.RawMessage `json:"signal,omitempty"`
}

type ChatMessaging struct {
	ToIds       [][]byte `json:"toIds"`
	RoomId      []byte   `json:"roomId,omitempty"`
	FromId      []byte   `json:"fromId"`
	ContentType string   `json:"contentType"`
	Content     string   `json:"content"`
}
