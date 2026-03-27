package dto

import "github.com/google/uuid"

type SendLikeRequest struct {
	ToId uuid.UUID `json:"toId"`
}
