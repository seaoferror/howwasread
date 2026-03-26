package service

import "errors"

var (
	ErrSignInWithApple          = errors.New("fail to apple sign in")
	ErrCheckEmail               = errors.New("fail to check email usability")
	ErrSignUpWithEmail          = errors.New("fail to sign up with email")
	ErrLoginWithEmail           = errors.New("email or password is(are) incorrect")
	ErrSendEmailOTP             = errors.New("fail to send email OTP")
	ErrVerifyEmailOTP           = errors.New("fail to verify email OTP")
	ErrPhoneNumberAlreadyLinked = errors.New("this phone number already linked with other account")
	ErrSendSMSOTP               = errors.New("fail to send SMS OTP")
	ErrVerifySMSOTP             = errors.New("fail to verify SMS OTP")
	ErrGenerateToken            = errors.New("fail to generate new access token")

	ErrInternalServer = errors.New("something went wrong")
)
