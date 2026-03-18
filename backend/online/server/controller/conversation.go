package controller

import (
	"backend/online/server/dto"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func conversationRouter(c *Controller) {
	c.Router(POST, "/online/conversation/create", c.createConversation)
	c.Router(GET, "/online/conversation/list", c.getConversations)
	c.Router(GET, "/online/conversation/join", c.joinConversation)
	c.Router(GET, "/online/conversation", c.getConversation)
}

func (c *Controller) createConversation(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateConversationRequest
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse userId from X-User-Id header",
			"err", err,
			"memberIdRaw", memberIdRaw,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Info("incorrect body",
			"err", err,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect request body"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}

	length, err := time.ParseDuration(req.Length)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		slog.Error("fail to parse duration from rawLength",
			"err", err,
			"req.Length", req.Length,
		)
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
		w.WriteHeader(getStatusCode(err))
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err,
		)
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
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	pageRaw := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		slog.Info("incorrect query param for page",
			"err", err,
			"pageRaw", pageRaw)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("wrong cursor"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	if page < 1 {
		page = 1
	}
	result, err := c.service.GetConversations(r.Context(), memberId, page)
	if err != nil {
		w.WriteHeader(getStatusCode(err))
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err)
		}
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
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	conversationIdRaw := r.URL.Query().Get("id")
	conversationId, err := bson.ObjectIDFromHex(conversationIdRaw)
	if err != nil {
		slog.Error("fail to parse conversation object id from raw string",
			"conversationIdRaw", conversationIdRaw)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect param"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	result, err := c.service.GetConversation(r.Context(), conversationId, memberId)
	if err != nil {
		w.WriteHeader(getStatusCode(err))
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}
