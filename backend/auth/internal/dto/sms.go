package dto

type SendSMSOTPRequest struct {
	SessionId   *string `json:"sessionId"` //nullable
	PhoneNumber string  `json:"phoneNumber"`
}
type VerifySMSOTPRequest struct {
	SessionId      *string `json:"sessionId"` //nullable
	VerificationId string  `json:"verificationId"`
	OTP            string  `json:"otp"`
}

type VerifySMSOTPResponse struct {
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	AccessToken         string `json:"accessToken"`
}
