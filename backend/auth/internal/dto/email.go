package dto

type CheckEmailRequest struct {
	Email string `json:"email"`
}
type SignInWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginWithEmailResponse struct {
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	EmailVerified       bool   `json:"emailVerified"`
	VerificationId      string `json:"verificationId"`
	SessionId           string `json:"sessionId"`
	AccessToken         string `json:"accessToken"`
}

type VerifyEmailOTPRequest struct {
	VerificationId string `json:"verificationId"`
	OTP            string `json:"otp"`
}
type VerifyEmailOTPResponse struct {
	EmailVerified bool   `json:"emailVerified"`
	SessionId     string `json:"session_id"`
}

type SignInWithAppleRequest struct {
	User          string  `json:"user"`
	Email         *string `json:"email"` //nullable
	IdentityToken string  `json:"identityToken"`
	Nonce         string  `json:"nonce"`
}

type SignInWithAppleResponse struct {
	PhoneNumberVerified bool   `json:"phoneNumberVerified"`
	SessionId           string `json:"sessionId"`
	AccessToken         string `json:"accessToken"`
}
