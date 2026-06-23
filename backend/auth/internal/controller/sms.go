package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"backend/auth/internal/constant"
	"backend/auth/internal/dto"
)

func smsRouter(n *Controller) {
	n.Router(POST, "/auth/sms/otp/send", n.sendSMSOTP)
	n.Router(POST, "/auth/sms/otp/verify", n.verifySMSOTP)
}

func (c *Controller) sendSMSOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.SendSMSOTPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}
	result, err := c.service.SendSMSOTP(req.SessionId, req.PhoneNumber)
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

func (c *Controller) verifySMSOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifySMSOTPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}
	result, rt, err := c.service.VerifySMSOTP(req.SessionId, req.VerificationId, req.OTP)
	if err != nil {
		handleError(w, err)
		return
	}
	if rt != "" {
		http.SetCookie(w, &http.Cookie{Name: "refresh_token",
			Value:    rt,
			Expires:  time.Now().Add(constant.RefreshTokenTTL * time.Second),
			HttpOnly: true,
			Secure:   true,
		})
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}
