package repository

import (
	"backend/auth/internal/constant"
	"context"
	"errors"
	"log/slog"

	"github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) SaveEmailLoginInfo(id gocql.UUID, email, password string) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query(
			`INSERT INTO member_by_email (
                             email_verified, phone_number_verified, id, email, password, role
                             ) VALUES (?, ?, ?, ?, ?, ?);`,
			false, false, id, email, password, constant.RoleUser).
		Query(
			`INSERT INTO member_by_id (
                          email_verified, phone_number_verified, id, email, role
                          ) VALUES (?, ?, ?, ?, ?)`,
			false, false, id, email, constant.RoleUser).
		Exec()
	if err != nil {
		slog.Error("fail to save member",
			"err", err,
			"id", id.String(),
		)
		return err
	}
	return nil
}

func (r *Repository) VerifiedEmailExists(ctx context.Context, email string) (bool, error) {
	var emailVerified bool
	err := r.session.Query(
		`SELECT email_verified FROM member_by_email WHERE email = ?`,
		email,
	).ScanContext(ctx, &emailVerified)
	if errors.Is(gocql.ErrNotFound, err) {
		return false, nil
	}
	if err != nil {
		slog.Error("fail to check email existence",
			"err", err,
			"email", email,
		)
		return true, err
	}
	return emailVerified, nil
}

func (r *Repository) FindLoginInfoByEmail(email string) (emailVerified, phoneNumberVerified bool, id gocql.UUID, password, role string, err error) {
	err = r.session.Query(
		`SELECT email_verified, phone_number_verified, id, password, role FROM member_by_email WHERE email = ?`,
		email,
	).Scan(&emailVerified, &phoneNumberVerified, &id, &password, &role)
	if err != nil {
		slog.Info("fail to find by email",
			"err", err,
			"email", email,
		)
		return false, false, gocql.UUID{}, "", "", err
	}
	return emailVerified, phoneNumberVerified, id, password, role, nil
}

func (r *Repository) SaveEmailAndOtpByVerificationId(verificationId gocql.UUID, email, otp string) error {
	err := r.session.Query(
		"INSERT INTO member_by_verification_id (verification_id, email, otp) VALUES (?, ?, ?) USING TTL ?",
		verificationId, email, otp, constant.AuthIdTTL,
	).Exec()
	if err != nil {
		slog.Error("fail to save email otp by verificationId",
			"err", err,
		)
		return err
	}
	return nil
}

func (r *Repository) FindEmailAndOTPByVerificationId(verificationId gocql.UUID) (email string, otp string, err error) {
	err = r.session.Query(
		"SELECT email, otp FROM member_by_verification_id WHERE verification_id = ?",
		verificationId,
	).Scan(&email, &otp)
	if err != nil {
		slog.Info("fail to select email and otp by verification_id",
			"err", err,
			"verificationId", verificationId.String(),
		)
		return "", "", err
	}
	return email, otp, nil
}

func (r *Repository) MarkEmailVerified(email string) error {
	var id gocql.UUID
	err := r.session.Query(
		"SELECT id FROM member_by_email WHERE email = ?",
		email,
	).Scan(&id)
	if err != nil {
		slog.Error("fail to select id by email",
			"err", err,
			"email", email,
		)
		return err
	}
	err = r.session.Query(
		"UPDATE member_by_id SET email_verified = ? WHERE id = ?",
		true, id,
	).Exec()
	if err != nil {
		slog.Error("fail to update email_verified at member_by_id",
			"err", err,
			"id", id.String(),
			"email", email,
		)
		return err
	}
	err = r.session.Query(
		"UPDATE member_by_email SET email_verified = ? WHERE email = ?",
		true, email,
	).Exec()
	if err != nil {
		slog.Error("fail to update email_verified at member_by_email",
			"err", err,
			"id", id.String(),
			"email", email,
		)
		return err
	}
	return nil
}

func (r *Repository) SaveEmailBySessionId(sessionId gocql.UUID, email string) error {
	err := r.session.Query(
		"INSERT INTO member_by_session_id (session_id, email) VALUES (?, ?) USING TTL ?",
		sessionId, email, constant.AuthIdTTL,
	).Exec()
	if err != nil {
		msg := "fail to insert email by session_id at email_by_session_id"
		slog.Error(msg,
			"err", err,
			"sessionId", sessionId.String(),
			"email", email,
		)
	}
	return nil
}

func (r *Repository) FindEmailBySessionId(sessionId gocql.UUID) (email string, err error) {
	err = r.session.Query(
		"SELECT email FROM member_by_session_id WHERE session_id = ?",
		sessionId,
	).Scan(&email)
	if err != nil {
		slog.Info("fail to find email by sessionId, might be expired sessionId",
			"err", err,
			"sessionId", sessionId,
		)
	}
	return email, nil
}
