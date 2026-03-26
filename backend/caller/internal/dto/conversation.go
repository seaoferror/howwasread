package dto

import "encoding/json"

type ConversationSignalMessage struct {
	FromId string          `json:"fromId"`
	ToId   string          `json:"toId"`
	Signal json.RawMessage `json:"signal,omitempty"`
}
