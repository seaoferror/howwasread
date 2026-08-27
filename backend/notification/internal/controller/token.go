package controller

import (
	"backend/notification/internal/dto"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func registerRouter(c *Controller) {
	c.Router(POST, "/notification/register", c.registerNotificationInfo)
}

func (c *Controller) registerNotificationInfo(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id",
			"err", err,
			"memberIdRaw", memberIdRaw)
		handleError(w, errors.New("incorrect body"))
	}
	var req dto.RegisterNotificationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("fail to parse body",
			"err", err)
		handleError(w, errors.New("incorrect body"))
	}
	err = c.service.RegisterNotification(r.Context(), memberId, req.OS, req.DevicePushToken)
	if err != nil {
		handleError(w, err)
	}
	w.WriteHeader(http.StatusOK)
	slog.Info("200 OK device push token",
		"id", memberId,
		"os", req.OS,
		"devicePushToken", req.DevicePushToken)
}
