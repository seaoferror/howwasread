package controller

import (
	"backend/chat/internal/dto"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func infoRouter(c *Controller) {
	c.Router(GET, "/chat/profile/my", c.getMyProfile)
	c.Router(PUT, "/chat/profile/name", c.setName)
	c.Router(GET, "/chat/profile", c.getProfile)
	c.Router(GET, "/chat/room/info", c.getChatRoomInfo)
}

func (c *Controller) getChatRoomInfo(w http.ResponseWriter, r *http.Request) {
	roomId, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	result, err := c.service.GetChatRoomInfo(r.Context(), roomId)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}

func (c *Controller) getProfile(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	result, err := c.service.GetProfile(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}

func (c *Controller) getMyProfile(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
		)
		handleError(w, errors.New("fail to parse"))
		return
	}
	result, err := c.service.GetProfile(r.Context(), memberId)
	if err != nil {
		handleError(w, err)
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
	memberId, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
		)
		handleError(w, errors.New("fail to parse"))
		return
	}
	var req dto.SetNameRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Info("incorrect body",
			"err", err,
		)
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.SetName(r.Context(), memberId, req.Name)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
