package payload

import (
	"encoding/json"
	"log/slog"

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

type ConversationRequest struct {
	Id uuid.UUID `json:"id"`
}
