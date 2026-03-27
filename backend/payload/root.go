package payload

import (
	"encoding/json"
)

type ConversationSignal struct {
	FromId json.RawMessage `json:"fromId"`
	ToId   json.RawMessage `json:"toId"`
	Signal json.RawMessage `json:"signal,omitempty"`
}
