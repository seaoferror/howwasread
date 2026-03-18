package dto

import (
	"encoding/json"
	"time"
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
	By           string    `json:"by,omitempty"`
	Rule         string    `json:"rule,omitempty"`
	When         time.Time `json:"when"`
	Length       string    `json:"length"`
	Ongoing      bool      `json:"ongoing"`
	IsModerator  bool      `json:"isModerator"`
	IsRegistrant bool      `json:"isRegistrant"`
}

type ConversationSignalResponse struct {
	FromIds []string        `json:"fromIds"`
	Signal  json.RawMessage `json:"signal,omitempty"`
}

type ConversationSignalRequest struct {
	ToIds  []string        `json:"toIds"`
	Signal json.RawMessage `json:"signal"`
}

type ConversationSignalMessage struct {
	FromId string          `json:"fromId"`
	ToId   string          `json:"toId"`
	Signal json.RawMessage `json:"signal,omitempty"`
}

type GetConversationResponse struct {
	Id         string    `json:"id"`
	Novel      string    `json:"novel,omitempty"`
	ShortStory string    `json:"shortStory,omitempty"`
	Poem       string    `json:"poem,omitempty"`
	Play       string    `json:"play,omitempty"`
	Film       string    `json:"film,omitempty"`
	By         string    `json:"by,omitempty"`
	Rule       string    `json:"rule,omitempty"`
	When       time.Time `json:"when"`
	Length     string    `json:"length"`

	IsModerator bool `json:"isModerator"`
}
