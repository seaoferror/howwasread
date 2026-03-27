package dto

import "github.com/google/uuid"

type SetNameRequest struct {
	Name string `json:"name"`
}

type GetMyProfileResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
