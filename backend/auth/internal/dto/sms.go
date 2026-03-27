package dto

import "github.com/google/uuid"

type SendSMSOTPRequest struct {
	SessionId   uuid.UUID `json:"sessionId"` //nullable
	PhoneNumber string    `json:"phoneNumber"`
}
type VerifySMSOTPRequest struct {
	SessionId      uuid.UUID `json:"sessionId"` //nullable
	VerificationId uuid.UUID `json:"verificationId"`
	OTP            string    `json:"otp"`
}

type VerifySMSOTPResponse struct {
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	AccessToken         string `json:"accessToken"`
}
