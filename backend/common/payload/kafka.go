package payload

import (
	"encoding/json"

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
	NotificationId uint8    `json:"notificationId,omitempty"`
	Id             []byte   `json:"id"`
	ToIds          [][]byte `json:"toIds"`
	RoomId         []byte   `json:"roomId"`
	FromId         []byte   `json:"fromId"`
	ContentType    string   `json:"contentType"`
	Contents       []string `json:"contents"`
}

type NotificationMessage struct {
	TokenMap map[string]uuid.UUID `json:"tokenMap"`
	Title    string               `json:"title,omitempty"`
	SubTitle string               `json:"subtitle"`
	Text     string               `json:"text"`
	ImageURL string               `json:"imageURL,omitempty"`
}

type NotificationScheduling struct {
	ScheduledTime int64          `json:"scheduledTime,omitempty"`
	PartitionId   uuid.UUID      `json:"partitionId"`
	KeyId         uuid.UUID      `json:"keyId"`
	Contents      map[int]string `json:"contents"`
	Type          string         `json:"type,omitempty"`
}

type NotificationScheduled struct {
	PartitionId    uuid.UUID                    `json:"partitionId"`
	PartitionType  string                       `json:"partitionType"`
	ScheduledTime  int64                        `json:"scheduledTime"`
	SharedContents map[int]string               `json:"sharedContents"`
	Notifications  map[uuid.UUID]map[int]string `json:"notifications"`
}

type RetryEvent struct {
	PartitionId uuid.UUID       `json:"partitionId"`
	Reason      string          `json:"reason"`
	Backoff     int64           `json:"backoff,omitempty"`
	Multiplier  int             `json:"multiplier,omitempty"`
	Cap         int             `json:"cap,omitempty"`
	MaxFailure  int64           `json:"maxFailure,omitempty"`
	Topic       string          `json:"topic,omitempty"`
	Type        string          `json:"type,omitempty"`
	Value       json.RawMessage `json:"value"`
}
