package payload

import "encoding/json"

type OnlineConversationSignal struct {
	FromId []byte          `json:"fromId"`
	ToIds  [][]byte        `json:"toIds"`
	Signal json.RawMessage `json:"signal,omitempty"`
}

type ChatMessaging struct {
	Id          []byte `json:"id"`
	FromId      []byte `json:"fromId"`
	ToIdType    string `json:"toIdType"`
	ToId        []byte `json:"toId"`
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}
