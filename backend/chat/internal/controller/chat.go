package controller

import (
	"backend/chat/internal/dto"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func chatRouter(c *Controller) {
	c.Router(POST, "/chat/like", c.sendLike)
}

func (c *Controller) sendLike(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		handleParseError(w)
		return
	}
	var req dto.SendLikeRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleParseError(w)
		return
	}
	toId, err := uuid.Parse(req.ToId)
	if err != nil {
		handleParseError(w)
		return
	}
	err = c.service.SendLike(r.Context(), memberId, toId)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
