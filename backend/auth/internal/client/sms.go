package client

import (
	"errors"
	"log/slog"
	"os"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// SMSClient sends and checks one time passwords by text message
type SMSClient interface {
	// SendOTP texts an OTP to phone and returns the number it was sent to, as normalized by twilio
	SendOTP(phone string) (string, error)
	// CheckOTP reports whether code is the OTP sent to phone
	CheckOTP(phone, code string) (bool, error)
}

type smsClient struct {
	verify     *verify.ApiService
	serviceSid string
}

func NewSMSClient() SMSClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   os.Getenv("TWILIO_API_KEY"),
		Password:   os.Getenv("TWILIO_API_SECRET"),
		AccountSid: os.Getenv("TWILIO_ACCOUNT_SID"),
	})
	return &smsClient{verify: client.VerifyV2, serviceSid: os.Getenv("TWILIO_SERVICE_SID")}
}

func (s *smsClient) SendOTP(phone string) (string, error) {
	params := &verify.CreateVerificationParams{}
	params.SetTo(phone)
	params.SetChannel("sms")
	resp, err := s.verify.CreateVerification(s.serviceSid, params)
	if err != nil {
		return "", err
	}
	if resp.To == nil {
		return "", errors.New("twilio returned no recipient")
	}
	if resp.Status != nil {
		slog.Info("success to send sms otp code", "phoneNumber", *resp.To, "status", *resp.Status)
	}
	return *resp.To, nil
}

func (s *smsClient) CheckOTP(phone, code string) (bool, error) {
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(code)
	resp, err := s.verify.CreateVerificationCheck(s.serviceSid, params)
	if err != nil {
		return false, err
	}
	if resp.Status == nil {
		return false, errors.New("twilio returned no status")
	}
	return *resp.Status == "approved", nil
}
