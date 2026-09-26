package client

import (
	"net/smtp"
	"os"
)

// SMTPClient sends emails from the service account
type SMTPClient interface {
	// SendOTP mails the email verification code to to
	SendOTP(to, otp string) error
}

type smtpClient struct{}

func NewSMTPClient() SMTPClient {
	return smtpClient{}
}

func (smtpClient) SendOTP(to, otp string) error {
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	message := "Subject: Howwasread : Verify your email\n" + headers + "\n\n" + otp + "\ncode is valid for 5 minutes"

	from := os.Getenv("FROM_EMAIL")
	auth := smtp.PlainAuth("", from, os.Getenv("FROM_EMAIL_PASSWORD"), os.Getenv("FROM_EMAIL_SMTP"))
	return smtp.SendMail(os.Getenv("SMTP_ADDR"), auth, from, []string{to}, []byte(message))
}
