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

type ConversationFeedResponse struct {
	Id           string    `json:"id"`
	Novel        string    `json:"novel,omitempty"`
	ShortStory   string    `json:"shortStory,omitempty"`
	Poem         string    `json:"poem,omitempty"`
	Play         string    `json:"play,omitempty"`
	Film         string    `json:"film,omitempty"`
	By           string    `json:"by"`
	Rule         string    `json:"rule,omitempty"`
	Capacity     int       `json:"capacity"`
	When         time.Time `json:"when"`
	Length       string    `json:"length"`
	Ongoing      bool      `json:"ongoing"`
	IsModerator  bool      `json:"isModerator"`
	IsRegistrant bool      `json:"isRegistrant"`
}

type ConversationSignalResponse struct {
	FromIds []uuid.UUID     `json:"fromIds"`
	Signal  json.RawMessage `json:"signal,omitempty"`
}

type ConversationSignalRequest struct {
	ToIds  []uuid.UUID     `json:"toIds"`
	Signal json.RawMessage `json:"signal"`
}

type GetConversationResponse struct {
	Id         string    `json:"id"`
	Novel      string    `json:"novel,omitempty"`
	ShortStory string    `json:"shortStory,omitempty"`
	Poem       string    `json:"poem,omitempty"`
	Play       string    `json:"play,omitempty"`
	Film       string    `json:"film,omitempty"`
	By         string    `json:"by"`
	Rule       string    `json:"rule,omitempty"`
	When       time.Time `json:"when"`
	Length     string    `json:"length"`

	IsModerator bool `json:"isModerator"`
}

type BanParticipantRequest struct {
	ConversationId string    `json:"conversationId"`
	BanId          uuid.UUID `json:"banId"`
}
