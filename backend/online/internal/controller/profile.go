package controller

import (
	"backend/online/internal/dto"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func profileRouter(c *Controller) {
	c.Router(GET, "/profile", c.getMyProfile)
	c.Router(PUT, "/profile/name", c.setName)
}

func (c *Controller) getMyProfile(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
			"memberIdRaw", memberIdRaw,
		)
		handleParseError(w)
		return
	}
	result, err := c.service.GetMyProfile(r.Context(), memberId)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}

func (c *Controller) setName(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
			"memberIdRaw", memberIdRaw,
		)
		handleParseError(w)
		return
	}
	var req dto.SetNameRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Info("incorrect body",
			"err", err,
		)
		handleParseError(w)
		return
	}
	err = c.service.SetName(r.Context(), memberId, req.Name)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
