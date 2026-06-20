package repository

import (
	"backend/auth/internal/constant"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) CheckNonce(nonce string) (bool, error) {
	var c int64
	err := r.session.Query(
		`SELECT COUNT(1) FROM nonce WHERE nonce = ?`,
		nonce,
	).Scan(&c)
	if c == 0 {
		return false, nil
	}
	if err != nil {
		slog.Error("fail to check nonce",
			"err", err,
			"nonce", nonce,
		)
		return false, err
	}
	return true, nil
}

func (r *Repository) SaveNonce(nonce string) error {
	err := r.session.Query(
		`INSERT INTO nonce (nonce) VALUES (?)`, nonce).Exec()
	if err != nil {
		slog.Error("fail to insert nonce",
			"err", err,
			"nonce", nonce)
		return err
	}
	return nil
}

func (r *Repository) SaveAppleSignInInfo(id gocql.UUID, appleSignInUser, email string, phoneNumberVerified, emailVerified bool) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query(
			`INSERT INTO member_by_id (
                          email_verified, phone_number_verified, id, apple_sign_in_user, email, role
                          ) VALUES (?, ?, ?, ?, ?, ?)`,
			true, phoneNumberVerified, id, appleSignInUser, email, constant.RoleUser).
		Query(
			`INSERT INTO member_by_apple_sign_in_user (
                                          id, apple_sign_in_user, email, role
                                          ) VALUES (?, ?, ?, ?)`,
			id, appleSignInUser, email, constant.RoleUser).
		Exec()
	if err != nil {
		slog.Error("fail to save apple sign in info",
			"err", err,
			"id", id.String(),
		)
		return err
	}
	if emailVerified {
		err = r.session.Query(
			`INSERT INTO member_by_email (
                             email_verified, phone_number_verified, id, apple_sign_in_user, email, role
                             ) VALUES (?, ?, ?, ?, ?, ?);`,
			emailVerified, phoneNumberVerified, id, appleSignInUser, email, constant.RoleUser).Exec()
		if err != nil {
			slog.Error("fail to save apple sign in info",
				"err", err,
				"id", id.String(),
			)
			return err
		}
		return nil
	}
	err = r.session.Query(
		`INSERT INTO member_by_email (
                             email_verified, password, phone_number_verified, id, apple_sign_in_user, email, role
                             ) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		true, nil, phoneNumberVerified, id, appleSignInUser, email, constant.RoleUser).Exec()
	return nil
}

func (r *Repository) FindAppleSignInInfoByUser(appleSignInUser string) (id gocql.UUID, email, role string, err error) {
	err = r.session.Query(
		"SELECT id, email, role FROM member_by_apple_sign_in_user WHERE apple_sign_in_user = ?",
		appleSignInUser,
	).Scan(&id, &email, &role)
	if err != nil {
		slog.Info("fail to find login info by apple sign in user",
			"err", err,
			"appleSignInUser", appleSignInUser,
		)
		return gocql.UUID{}, "", "", err
	}
	return id, email, role, nil
}
