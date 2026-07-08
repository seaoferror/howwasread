package controller

import (
	"backend/auth/internal/constant"
	"backend/auth/internal/dto"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func emailRouter(n *Controller) {
	n.Router(POST, "/auth/email/create", n.createMemberByEmail)
	n.Router(POST, "/auth/email/login", n.loginWithEmail)
	n.Router(POST, "/auth/email/otp/verify", n.verifyEmailOTP)
	n.Router(POST, "/auth/email/apple", n.signInWithApple)
	n.Router(POST, "/auth/email/google", n.signInWithGoogle)
}

func (c *Controller) createMemberByEmail(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInWithEmailRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}
	result, err := c.service.CreateMemberByEmail(r.Context(), req.Email, req.Password)
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

func (c *Controller) loginWithEmail(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInWithEmailRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}
	result, rt, err := c.service.LoginWithEmail(req.Email, req.Password)
	if err != nil {
		handleError(w, err)
		return
	}
	if rt != "" {
		http.SetCookie(w, &http.Cookie{Name: "refresh_token",
			Value:    rt,
			Expires:  time.Now().Add(constant.RefreshTokenTTL * time.Second),
			Path:     "/",
			Domain:   "",
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

func (c *Controller) verifyEmailOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyEmailOTPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}
	result, err := c.service.VerifyEmailOTP(req.OTP, req.VerificationId)
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

func (c *Controller) signInWithApple(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInWithAppleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}
	responseBody, rt, err := c.service.SignInWithApple(
		r.Context(),
		req.IdentityToken,
	)
	if err != nil {
		handleError(w, err)
		return
	}
	if rt != "" {
		http.SetCookie(w, &http.Cookie{Name: "refresh_token",
			Value:    rt,
			Expires:  time.Now().Add(constant.RefreshTokenTTL * time.Second),
			Path:     "/",
			Domain:   "",
			HttpOnly: true,
			Secure:   true,
		})
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}

func (c *Controller) signInWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInWithGoogleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}
	responseBody, rt, err := c.service.SignInWithGoogle(r.Context(), req.IdToken)
	if err != nil {
		handleError(w, err)
		return
	}
	if rt != "" {
		http.SetCookie(w, &http.Cookie{Name: "refresh_token",
			Value:    rt,
			Expires:  time.Now().Add(constant.RefreshTokenTTL * time.Second),
			Path:     "/",
			Domain:   "",
			HttpOnly: true,
			Secure:   true,
		})
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		slog.Error("fail to write response body",
			"err", err)
	}
}
