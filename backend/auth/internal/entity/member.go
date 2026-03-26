package entity

import (
	"time"

	"github.com/apache/cassandra-gocql-driver/v2"
)

type Member struct {
	//VerificationId is for otp verifying. By this, we check whether user
	//who try to verify otp is the same person who requested otp.
	//We also use this for get email, phone number, otp from db
	//for verification and mark member as verified after verification
	//by retrieved phone number and email
	VerificationId gocql.UUID `db:"verification_id"`
	//SessionId is for connection with previous verified email password login
	//Give this to user who success to verify their email and who need to verify
	//their sms(also give this to who success to log in and already verified their
	//email before, but who need to verify their sms)
	SessionId           gocql.UUID `db:"session_id"`
	OTP                 string     `db:"otp"`
	EmailVerified       bool       `db:"email_verified"`
	PhoneNumberVerified bool       `db:"phone_number_verified"`

	Id          gocql.UUID `db:"id"`
	PhoneNumber string     `db:"phone_number"`
	Email       string     `db:"email"`
	Password    string     `db:"password"`
	Name        string     `db:"name"`
	Role        string     `db:"role"`
	CreatedTime time.Time  `db:"created_time"`
	UpdatedTime time.Time  `db:"updated_time"`
	DeletedTime time.Time  `db:"deleted_time"`

	//Secret      string     `db:"secret"`
}
