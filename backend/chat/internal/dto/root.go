package dto

import "github.com/google/uuid"

type SendLikeRequest struct {
	ToId uuid.UUID `json:"toId"`
}

type SendMessagingRequest struct {
	ToIdType    string    `json:"toIdType"`
	ToId        uuid.UUID `json:"toId"`
	ContentType string    `json:"contentType"`
	Content     string    `json:"content"`
}

type MessagingResponse struct {
	Id          uuid.UUID `json:"id"`
	RoomId      uuid.UUID `json:"roomId,omitempty"`
	FromId      uuid.UUID `json:"fromId"`
	ContentType string    `json:"contentType"`
	Content     string    `json:"content"`
}
