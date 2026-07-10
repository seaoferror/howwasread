package dto

import "github.com/google/uuid"

type SendMessagingRequest struct {
	ToIdType    string    `json:"toIdType"`
	ToId        uuid.UUID `json:"toId"`
	ContentType string    `json:"contentType"`
	Contents    []string  `json:"contents"`
}

type MessagingResponse struct {
	Id          uuid.UUID `json:"id"`
	RoomId      uuid.UUID `json:"roomId"`
	FromId      uuid.UUID `json:"fromId"`
	ContentType string    `json:"contentType"`
	Contents    []string  `json:"contents"`
}

type GetChatRoomInfoResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SetNameRequest struct {
	Name string `json:"name"`
}

type GetProfileResponse struct {
	Id   uuid.UUID `json:"id,omitzero"`
	Name string    `json:"name,omitempty"`
}

type GeneratePresignedURLRequest struct {
	ContentType string `json:"contentType"`
	N           int    `json:"n"`
}

type GeneratePresignedURLResponse struct {
	Filename uuid.UUID         `json:"filename"`
	URL      string            `json:"url"`
	Fields   map[string]string `json:"fields"`
}

type BlockReport struct {
	Id string `json:"id"`
}
