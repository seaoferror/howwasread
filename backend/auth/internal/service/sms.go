package service

import (
	"backend/auth/internal/constant"
	"backend/auth/internal/dto"
	"errors"
	"log/slog"
	"os"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	verify "github.com/twilio/twilio-go/rest/verify/v2"

	_ "github.com/joho/godotenv/autoload"
)

func (s *Service) SendSMSOTP(sessionId *string, phoneNumber string) (map[string]string, error) {
	if sessionId != nil {
		email, err := s.repository.FindEmailByPhoneNumber(phoneNumber)
		if errors.Is(err, gocql.ErrNotFound) {
			err = nil
		}
		if err != nil {
			return nil, ErrSendSMSOTP
		}
		if email != "" {
			//TODO: ban this account by counting if he kept trying to use already linked phone number
			return nil, ErrPhoneNumberAlreadyLinked
		}
	}

	serviceSid := os.Getenv("TWILIO_SERVICE_SID")

	params := &verify.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := s.twilioClient.VerifyV2.CreateVerification(serviceSid, params)
	if err != nil {
		slog.Info("fail to send sms otp code",
			"err", err,
			"phoneNumber", phoneNumber,
		)
		return nil, ErrSendSMSOTP
	}
	if resp.Status != nil {
		slog.Info("success to send sms otp code",
			"phoneNumber", resp.To,
			"status", *resp.Status,
		)
	}
	vid, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to make random uuid for verification id")
		return nil, ErrInternalServer
	}
	err = s.repository.SavePhoneNumberByVerificationId(vid, *resp.To)
	if err != nil {
		return nil, ErrInternalServer
	}
	res := map[string]string{"verificationId": vid.String()}
	return res, nil
}

func (s *Service) VerifySMSOTP(sessionId *string, otp, verificationId string) (*dto.VerifySMSOTPResponse, string /*refreshToken*/, error) {
	var email string
	if sessionId != nil {
		sid, err := gocql.ParseUUID(*sessionId)
		if err != nil {
			slog.Info("fail to parse SessionId from OTPVerifyReq",
				"err", err,
				"SessionId", *sessionId,
			)
			return nil, "", ErrVerifySMSOTP
		}

		email, err = s.repository.FindEmailBySessionId(sid)
		if err != nil {
			return nil, "", ErrVerifySMSOTP
		}
	}

	vid, err := gocql.ParseUUID(verificationId)
	if err != nil {
		slog.Info("fail to parse verificationId from request",
			"err", err,
			"id", verificationId,
		)
		return nil, "", ErrVerifySMSOTP
	}
	phoneNumber, err := s.repository.FindPhoneNumberByVerificationId(vid)
	if err != nil {
		return nil, "", ErrVerifySMSOTP
	}

	if sessionId != nil {
		e, err := s.repository.FindEmailByPhoneNumber(phoneNumber)
		if errors.Is(err, gocql.ErrNotFound) {
			err = nil
		}
		if err != nil {
			return nil, "", ErrVerifySMSOTP
		}
		if e != "" {
			return nil, "", ErrVerifySMSOTP
		}
	}

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(otp)

	serviceSid := os.Getenv("TWILIO_SERVICE_SID")

	resp, err := s.twilioClient.VerifyV2.CreateVerificationCheck(serviceSid, params)
	if err != nil {
		slog.Error("fail to verify phone number otp",
			"err", err,
		)
		return nil, "", ErrVerifySMSOTP
	}
	if resp.Status == nil {
		slog.Error("status is nil pointer")
		return nil, "", ErrVerifySMSOTP
	}
	if *resp.Status != "approved" {
		slog.Info("otp is not correct",
			"otp", otp,
		)
		r := dto.VerifySMSOTPResponse{
			PhoneNumberVerified: false,
		}
		return &r, "", nil
	}

	if sessionId == nil {
		var idv7 uuid.UUID
		id, err := s.repository.FindIdByPhoneNumber(phoneNumber)
		if errors.Is(err, gocql.ErrNotFound) {
			err = nil
			idv7, err = uuid.NewV7()
			if err != nil {
				slog.Error("fail to create uuid v7 for phone number sign in user")
				return nil, "", ErrInternalServer
			}
			id = gocql.UUID(idv7)
		}
		if err != nil {
			slog.Error("fail to find id by phone number except for not found",
				"err", err,
			)
			return nil, "", ErrInternalServer
		}
		err = s.repository.SavePhoneNumberLoginInfo(phoneNumber, id)
		if err != nil {
			return nil, "", ErrInternalServer
		}
		jti, err := gocql.RandomUUID()
		if err != nil {
			slog.Error("fail to make random uuid for jti")
			return nil, "", ErrInternalServer
		}
		at, rt, err := s.createLoginTokens(id.String(), jti.String(), constant.RoleUser)
		if err != nil {
			return nil, "", ErrInternalServer
		}
		err = s.repository.SaveRefreshTokenJTIById(id, jti)
		if err != nil {
			return nil, "", ErrInternalServer
		}
		r := dto.VerifySMSOTPResponse{
			PhoneNumberVerified: true,
			AccessToken:         at,
		}
		if idv7 != uuid.Nil {
			err = s.kafkaProducer.PushMessage("auth.new_member_id", id.Bytes())
			if err != nil {
				//TODO: Retry logic with grpc or rest
				return &r, rt, nil
			}
		}

		return &r, rt, nil
	}

	_, _, id, password, role, err := s.repository.FindLoginInfoByEmail(email)
	if err != nil {
		slog.Warn("fail to find login info by email which is selected by sessionId",
			"err", err,
		)
		return nil, "", ErrInternalServer
	}

	oldAccountId, err := s.repository.FindIdByPhoneNumber(phoneNumber)
	if err == nil {
		err = s.repository.ReplaceAndLinkMemberWithOldAccount(id, oldAccountId, email, password, phoneNumber)
		id = oldAccountId
	}
	if errors.Is(err, gocql.ErrNotFound) {
		err = nil
		err = s.repository.LinkAndMarkVerifiedPhoneNumber(id, email, phoneNumber, role)
	}
	if err != nil {
		return nil, "", ErrInternalServer
	}

	jti, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to make random uuid for jti")
		return nil, "", ErrInternalServer
	}
	at, rt, err := s.createLoginTokens(id.String(), jti.String(), constant.RoleUser)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	err = s.repository.SaveRefreshTokenJTIById(id, jti)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	r := dto.VerifySMSOTPResponse{
		PhoneNumberVerified: true,
		AccessToken:         at,
	}
	err = s.kafkaProducer.PushMessage("auth.new_member_id", id.Bytes())
	if err != nil {
		//TODO: Retry logic with grpc or rest
		return &r, rt, nil
	}
	return &r, rt, nil
}
