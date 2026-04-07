package controller

import (
	"backend/notification/internal/dto"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func tokenRouter(c *Controller) {
	c.Router(POST, "/notification/device-push-token", c.registerDevicePushToken)
}

func (c *Controller) registerDevicePushToken(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id",
			"err", err,
			"memberIdRaw", memberIdRaw)
		handleError(w, errors.New("incorrect body"))
	}
	var req dto.SetDevicePushTokenRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("fail to parse body",
			"err", err)
		handleError(w, errors.New("incorrect body"))
	}
	err = c.service.SetDevicePushToken(r.Context(), memberId, req.DeviceId, req.OS, req.DevicePushToken)
	if err != nil {
		handleError(w, err)
	}
	w.WriteHeader(http.StatusOK)
	slog.Info("200 OK device push token",
		"id", memberId,
		"deviceId", req.DeviceId,
		"os", req.OS,
		"devicePushToken", req.DevicePushToken)
}
