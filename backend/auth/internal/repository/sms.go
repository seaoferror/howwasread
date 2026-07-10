package repository

import (
	"backend/auth/internal/constant"
	"log/slog"

	"github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) SavePhoneNumberByVerificationId(verificationId gocql.UUID, phoneNumber string) error {
	err := r.session.Query("INSERT INTO member_by_verification_id (phone_number, verification_id) values (?,?) USING TTL ?",
		phoneNumber, verificationId, constant.OtpTTL,
	).Exec()
	if err != nil {
		slog.Error("fail to insert phone number with id",
			"err", err,
			"verificationId", verificationId,
			"phoneNumber", phoneNumber,
		)
		return err
	}
	return nil
}

func (r *Repository) FindPhoneNumberByVerificationId(verificationId gocql.UUID) (phoneNumber string, err error) {
	err = r.session.Query(
		"SELECT phone_number FROM member_by_verification_id WHERE verification_id = ?",
		verificationId,
	).Scan(&phoneNumber)
	if err != nil {
		slog.Info("fail to find phone_number by verification_id",
			"err", err,
			"verificationId", verificationId,
		)
		return "", err
	}
	return phoneNumber, nil
}

func (r *Repository) SavePhoneNumberLoginInfo(phoneNumber string, id gocql.UUID) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query(
			"INSERT INTO member_by_phone_number (phone_number_verified, id, phone_number, role) VALUES (?, ?, ?, ?)",
			true, id, phoneNumber, constant.RoleUser,
		).
		Query(
			"INSERT INTO member_by_id (phone_number_verified, id, phone_number, role) VALUES (?, ?, ?, ?)",
			true, id, phoneNumber, constant.RoleUser,
		).Exec()
	if err != nil {
		slog.Error("fail to insert member at member_by_id",
			"err", err,
			"phoneNumber", phoneNumber,
		)
		return err
	}
	return nil
}

func (r *Repository) LinkAndMarkVerifiedPhoneNumber(id gocql.UUID, email, phoneNumber, role string) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query("UPDATE member_by_email SET phone_number_verified = ?, phone_number = ? WHERE email = ?",
			true, phoneNumber, email).
		Query("UPDATE member_by_id SET phone_number_verified = ?, phone_number = ? WHERE id = ?",
			true, phoneNumber, id).
		Query("INSERT INTO member_by_phone_number (phone_number_verified, id, email, phone_number, role) VALUES (?, ?, ?, ?, ?)",
			true, id, email, phoneNumber, role).
		Exec()
	if err != nil {
		slog.Error("fail to set phone_number",
			"err", err,
			"id", id,
			"email", email,
			"phoneNumber", phoneNumber,
		)
		return err
	}
	return nil
}

func (r *Repository) FindIdByPhoneNumber(phoneNumber string) (id gocql.UUID, err error) {
	err = r.session.Query(
		"SELECT id FROM member_by_phone_number WHERE phone_number = ?",
		phoneNumber,
	).Scan(&id)
	if err != nil {
		slog.Info("fail to find id by phone number")
		return gocql.UUID{}, err
	}
	return id, nil
}

func (r *Repository) FindEmailByPhoneNumber(phoneNumber string) (email string, err error) {
	err = r.session.Query(
		"SELECT email FROM member_by_phone_number WHERE phone_number = ?",
		phoneNumber,
	).Scan(&email)
	if err != nil {
		slog.Info("fail to find email by phone number",
			"err", err,
			"phoneNumber", phoneNumber,
		)
		return "", err
	}
	return email, nil
}

func (r *Repository) ReplaceAndLinkMemberWithOldAccount(newId, oldAccountId gocql.UUID, email, phoneNumber string) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query("DELETE FROM member_by_id WHERE id = ?",
			newId).
		Query("UPDATE member_by_id SET email = ?, email_verified = ? WHERE id = ?",
			email, true, oldAccountId,
		).
		Query("UPDATE member_by_email SET id = ?, phone_number = ?, phone_number_verified = ? WHERE email = ?",
			oldAccountId, phoneNumber, true, email,
		).
		Query("UPDATE member_by_phone_number SET email = ? WHERE phone_number = ?",
			email, phoneNumber).
		Exec()
	if err != nil {
		slog.Error("fail to replace and link member with old phone number account", "err", err)
		return err
	}
	return nil
}

func (r *Repository) WasBanned(phoneNumber string) error {
	var p string
	err := r.session.Query(`SELECT phone_number FROM banned_phone_number WHERE phone_number = ?`, phoneNumber).Scan(&p)
	if err != nil {
		slog.Info("fail to check whether phone number was banned",
			"err", err)
		return err
	}
	return nil
}
