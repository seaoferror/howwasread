package service

import (
	"backend/auth/internal/constant"
	"backend/auth/internal/dto"
	"errors"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"

	_ "github.com/joho/godotenv/autoload"
)

func (s *service) SendSMSOTP(sessionId uuid.UUID, phoneNumber string) (map[string]uuid.UUID, error) {
	err := s.repository.WasBanned(phoneNumber)
	if err == nil {
		return nil, errors.New("this phone number is not usable")
	}
	if errors.Is(err, gocql.ErrNotFound) {
		err = nil
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	if sessionId == uuid.Nil {
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

	to, err := s.smsClient.SendOTP(phoneNumber)
	if err != nil {
		slog.Info("fail to send sms otp code",
			"err", err,
			"phoneNumber", phoneNumber,
		)
		return nil, ErrSendSMSOTP
	}
	vid, err := gocql.RandomUUID()
	if err != nil {
		slog.Error("fail to make random uuid for verification id")
		return nil, ErrInternalServer
	}
	err = s.repository.SavePhoneNumberByVerificationId(vid, to)
	if err != nil {
		return nil, ErrInternalServer
	}
	res := map[string]uuid.UUID{"verificationId": uuid.UUID(vid)}
	return res, nil
}

func (s *service) VerifySMSOTP(sessionId uuid.UUID, verificationId uuid.UUID, otp string) (*dto.VerifySMSOTPResponse, string, error) {
	var email string
	var err error
	if sessionId != uuid.Nil {
		email, err = s.repository.FindEmailBySessionId(gocql.UUID(sessionId))
		if err != nil {
			return nil, "", ErrVerifySMSOTP
		}
	}
	phoneNumber, err := s.repository.FindPhoneNumberByVerificationId(gocql.UUID(verificationId))
	if err != nil {
		return nil, "", ErrVerifySMSOTP
	}
	if sessionId != uuid.Nil {
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

	approved, err := s.smsClient.CheckOTP(phoneNumber, otp)
	if err != nil {
		slog.Error("fail to verify phone number otp",
			"err", err,
		)
		return nil, "", ErrVerifySMSOTP
	}
	if !approved {
		slog.Info("otp is not correct", "otp", otp)
		return nil, "", ErrVerifySMSOTP
	}

	if sessionId == uuid.Nil {
		var idv7 uuid.UUID
		id, err1 := s.repository.FindIdByPhoneNumber(phoneNumber)
		if errors.Is(err1, gocql.ErrNotFound) {
			err1 = nil
			idv7, err1 = uuid.NewV7()
			if err1 != nil {
				slog.Error("fail to create uuid v7 for phone number sign in user")
				return nil, "", ErrInternalServer
			}
			id = gocql.UUID(idv7)
		}
		if err1 != nil {
			slog.Error("fail to find id by phone number except for not found",
				"err", err1,
			)
			return nil, "", ErrInternalServer
		}
		err1 = s.repository.SavePhoneNumberLoginInfo(phoneNumber, id)
		if err1 != nil {
			return nil, "", ErrInternalServer
		}
		jti, err1 := gocql.RandomUUID()
		if err1 != nil {
			slog.Error("fail to make random uuid for jti")
			return nil, "", ErrInternalServer
		}
		at, rt, err1 := s.createLoginTokens(id.String(), jti.String(), constant.RoleUser)
		if err1 != nil {
			return nil, "", ErrInternalServer
		}
		err1 = s.repository.SaveRefreshTokenJTIById(id, jti)
		if err1 != nil {
			return nil, "", ErrInternalServer
		}
		r := dto.VerifySMSOTPResponse{
			AccessToken: at,
		}
		if idv7 != uuid.Nil {
			err1 = s.repository.SaveProfileId(id)
			if err1 != nil {
				return nil, "", ErrInternalServer
			}
		}
		return &r, rt, nil
	}

	_, _, id, _, role, err := s.repository.FindLoginInfoByEmail(email)
	if err != nil {
		slog.Warn("fail to find login info by email which is selected by sessionId",
			"err", err,
		)
		return nil, "", ErrInternalServer
	}

	oldAccountId, err := s.repository.FindIdByPhoneNumber(phoneNumber)
	if err == nil {
		err = s.repository.ReplaceAndLinkMemberWithOldAccount(id, oldAccountId, email, phoneNumber)
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
		AccessToken: at,
	}
	err = s.repository.SaveProfileId(id)
	if err != nil {
		return nil, "", ErrInternalServer
	}
	return &r, rt, nil
}
