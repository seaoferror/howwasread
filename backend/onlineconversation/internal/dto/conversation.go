package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateConversationRequest struct {
	Novel      string    `json:"novel"`
	ShortStory string    `json:"short_story"`
	Poem       string    `json:"poem"`
	Play       string    `json:"play"`
	Film       string    `json:"film"`
	By         string    `json:"by"`
	Rule       string    `json:"rule"`
	Capacity   int       `json:"capacity"`
	When       time.Time `json:"when"`
	Length     string    `json:"length"`
}

type OnlineConversationDetailResponse struct {
	Id           uuid.UUID   `json:"id"`
	Novel        string      `json:"novel,omitempty"`
	ShortStory   string      `json:"shortStory,omitempty"`
	Poem         string      `json:"poem,omitempty"`
	Play         string      `json:"play,omitempty"`
	Film         string      `json:"film,omitempty"`
	By           string      `json:"by"`
	Rule         string      `json:"rule,omitempty"`
	Capacity     int         `json:"capacity"`
	When         time.Time   `json:"when"`
	Length       string      `json:"length"`
	CanEnter     bool        `json:"canEnter"`
	IsRegistrant bool        `json:"isRegistrant"`
	ModeratorIds []uuid.UUID `json:"moderatorIds"`
}

type OnlineConversationDocument struct {
	Id         uuid.UUID `json:"id"`
	Novel      string    `json:"novel"`
	ShortStory string    `json:"short_story"`
	Poem       string    `json:"poem"`
	Play       string    `json:"play"`
	Film       string    `json:"film"`
	WrittenBy  string    `json:"writtenBy"`
	Time       time.Time `json:"time"`
}

type ConversationSignalResponse struct {
	FromIds []uuid.UUID     `json:"fromIds"`
	Signal  json.RawMessage `json:"signal,omitempty"`
}

type ConversationSignalRequest struct {
	ToIds  []uuid.UUID     `json:"toIds"`
	Signal json.RawMessage `json:"signal"`
}

type OnlineConversationFeedResponse struct {
	Id         uuid.UUID `json:"id"`
	Novel      string    `json:"novel,omitempty"`
	ShortStory string    `json:"shortStory,omitempty"`
	Poem       string    `json:"poem,omitempty"`
	Play       string    `json:"play,omitempty"`
	Film       string    `json:"film,omitempty"`
	By         string    `json:"by"`
	When       time.Time `json:"when"`
}

type BanParticipantRequest struct {
	ConversationId uuid.UUID `json:"conversationId"`
	BanId          uuid.UUID `json:"banId"`
}

type ConversationRequest struct {
	Id uuid.UUID `json:"id"`
}
