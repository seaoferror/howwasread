package repository

import (
	"backend/auth/internal/constant"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) SaveAppleSignInInfo(id gocql.UUID, appleSignInUser, email string, phoneNumberVerified bool) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query(
			"INSERT INTO member_by_email (email_verified, phone_number_verified, id, apple_sign_in_user, email, role) VALUES (?, ?, ?, ?, ?, ?);",
			true, phoneNumberVerified, id, appleSignInUser, email, constant.RoleUser).
		Query(
			"INSERT INTO member_by_id (email_verified, phone_number_verified, id, apple_sign_in_user, email, role) VALUES (?, ?, ?, ?, ?, ?)",
			true, phoneNumberVerified, id, appleSignInUser, email, constant.RoleUser).
		Query(
			"INSERT INTO member_by_apple_sign_in_user (email_verified, id, apple_sign_in_user, email, role) VALUES (?, ?, ?, ?, ?)",
			true, id, appleSignInUser, email, constant.RoleUser).
		Exec()
	if err != nil {
		slog.Error("fail to save apple sign in info",
			"err", err,
			"id", id.String(),
		)
		return err
	}
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
