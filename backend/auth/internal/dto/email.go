package dto

import "github.com/google/uuid"

type CheckEmailRequest struct {
	Email string `json:"email"`
}
type SignInWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginWithEmailResponse struct {
	VerificationId uuid.UUID `json:"verificationId"`
	SessionId      uuid.UUID `json:"sessionId"`
	AccessToken    string    `json:"accessToken"`
}

type VerifyEmailOTPRequest struct {
	VerificationId uuid.UUID `json:"verificationId"`
	OTP            string    `json:"otp"`
}
type VerifyEmailOTPResponse struct {
	SessionId uuid.UUID `json:"session_id"`
}

type SignInWithAppleRequest struct {
	IsFirstSignIn bool   `json:"isFirstSignIn"`
	IdentityToken string `json:"identityToken"`
}

type SignInWithGoogleRequest struct {
	IdToken string `json:"idToken"`
}

type SignInWithThirdPartyResponse struct {
	SessionId   uuid.UUID `json:"sessionId"`
	AccessToken string    `json:"accessToken"`
}
