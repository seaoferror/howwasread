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
	"go.mongodb.org/mongo-driver/v2/bson"
)

func conversationRouter(c *Controller) {
	c.Router(POST, "/onlineconversation/conversation/create", c.createConversation)
	c.Router(GET, "/onlineconversation/conversation/list", c.getConversations)
	c.Router(GET, "/onlineconversation/conversation/join", c.joinConversation)
	c.Router(GET, "/onlineconversation/conversation", c.getConversation)
	c.Router(POST, "/onlineconversation/conversation/ban", c.banParticipant)
	c.Router(POST, "/onlineconversation/conversation/report", c.reportOnlineConversation)
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
		req.By,
		req.Rule,
		req.Capacity,
		req.When,
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
	pageRaw := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		slog.Info("incorrect query param for page",
			"err", err,
			"pageRaw", pageRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}
	if page < 1 {
		page = 1
	}
	result, err := c.service.GetConversations(r.Context(), memberId, page)
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

func (c *Controller) getConversation(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err,
			"memberIdRaw", memberIdRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}
	conversationIdRaw := r.URL.Query().Get("id")
	conversationId, err := bson.ObjectIDFromHex(conversationIdRaw)
	if err != nil {
		slog.Error("fail to parse conversation object id from raw string",
			"conversationIdRaw", conversationIdRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}
	result, err := c.service.GetConversation(r.Context(), conversationId, memberId)
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
	conversationId, err := bson.ObjectIDFromHex(req.ConversationId)
	if err != nil {
		slog.Error("fail to parse conversation id from raw string", "err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.BanParticipant(r.Context(), memberId, conversationId, req.BanId)
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
	var req dto.ReportConversationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	conversationId, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Error("fail to parse conversation id from raw string", "err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.ReportOnlineConversation(r.Context(), memberId, conversationId)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
