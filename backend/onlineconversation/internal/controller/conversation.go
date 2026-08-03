package controller

import (
	"backend/onlineconversation/internal/dto"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func conversationRouter(c *Controller) {
	c.Router(POST, "/onlineconversation/create", c.createConversation)
	c.Router(GET, "/onlineconversation/list", c.getConversations)
	c.Router(GET, "/onlineconversation/join", c.joinConversation)
	c.Router(GET, "/onlineconversation/detail", c.getConversationDetail)
	c.Router(POST, "/onlineconversation/ban", c.banParticipant)
	c.Router(POST, "/onlineconversation/report", c.reportOnlineConversation)
	c.Router(POST, "/onlineconversation/register", c.registerOnlineConversation)
	c.Router(POST, "/onlineconversation/deregister", c.deregisterOnlineConversation)
}

func (c *Controller) createConversation(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
			"memberIdRaw", memberIdRaw,
		)
		handleError(w, errors.New("fail to parse"))
		return
	}
	var req dto.CreateConversationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Info("incorrect body",
			"err", err,
		)
		handleError(w, errors.New("fail to parse"))
		return
	}

	length, err := time.ParseDuration(req.Length)
	if err != nil {
		slog.Error("fail to parse duration from rawLength",
			"err", err,
			"req.Length", req.Length,
		)
		handleError(w, errors.New("fail to parse"))
		return
	}

	result, err := c.service.CreateConversation(
		r.Context(),
		memberId,
		req.Novel,
		req.ShortStory,
		req.Poem,
		req.Play,
		req.Film,
		req.WrittenBy,
		req.Rule,
		req.Capacity,
		req.Time,
		length,
	)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body", "err", err)
	}
}

func (c *Controller) getConversations(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		slog.Info("incorrect query param for page",
			"err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	if page < 1 {
		page = 1
	}
	t, err := time.Parse(time.RFC3339, r.URL.Query().Get("time"))
	if err != nil {
		t = time.Now()
		//slog.Info("incorrect query param for time",
		//	"err", err)
		//handleError(w, errors.New("fail to parse"))
		//return
		//TODO: uncomment this after deploy app
	}
	result, err := c.service.GetConversations(r.Context(), page, t)
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

func (c *Controller) getConversationDetail(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err,
			"memberIdRaw", memberIdRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}
	conversationId, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		slog.Error("fail to parse conversation uuid from raw string", "err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	result, err := c.service.GetConversationDetail(r.Context(), conversationId, memberId)
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

func (c *Controller) banParticipant(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err,
			"memberIdRaw", memberIdRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}
	var req dto.BanParticipantRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.BanParticipant(r.Context(), memberId, req.ConversationId, req.BanId)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) reportOnlineConversation(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	var req dto.ConversationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.ReportOnlineConversation(r.Context(), memberId, req.Id)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) registerOnlineConversation(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	var req dto.ConversationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.RegisterOnlineConversation(r.Context(), memberId, req.Id)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) deregisterOnlineConversation(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	var req dto.ConversationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.DeregisterOnlineConversation(r.Context(), memberId, req.Id)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
