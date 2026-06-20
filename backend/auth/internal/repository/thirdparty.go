package repository

import (
	"backend/auth/internal/constant"
	"context"
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

func (r *Repository) SaveThirdPartySignInInfo(ctx context.Context, id gocql.UUID, thirdPartyId, email string, phoneNumberVerified, emailVerified bool, providerName string) error {
	err := r.session.Batch(gocql.LoggedBatch).
		Query(
			`INSERT INTO member_by_id (
                          email_verified, phone_number_verified, id, email, role
                          ) VALUES (?, ?, ?, ?, ?)`,
			true, phoneNumberVerified, id, email, constant.RoleUser).
		Query(
			`INSERT INTO member_by_third_party (
                                          id, third_party_id, email, role, provider_name
                                          ) VALUES (?, ?, ?, ?, ?)`,
			id, thirdPartyId, email, constant.RoleUser, providerName).
		ExecContext(ctx)
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
                             email_verified, phone_number_verified, id, email, role
                             ) VALUES (?, ?, ?, ?, ?);`,
			emailVerified, phoneNumberVerified, id, email, constant.RoleUser).ExecContext(ctx)
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
                             email_verified, password, phone_number_verified, id, email, role
                             ) VALUES (?, ?, ?, ?, ?, ?);`,
		true, nil, phoneNumberVerified, id, email, constant.RoleUser).ExecContext(ctx)
	return nil
}

func (r *Repository) FindThirdPartySignInInfo(ctx context.Context, thirdPartyId string, providerName string) (id gocql.UUID, email, role string, err error) {
	err = r.session.Query(
		"SELECT id, email, role FROM member_by_third_party WHERE provider_name = ? AND third_party_id = ?",
		providerName, thirdPartyId,
	).ScanContext(ctx, &id, &email, &role)
	if err != nil {
		slog.Info("fail to find login info by third party",
			"err", err,
			"thirdPartyId", thirdPartyId,
			"providerName", providerName)
		return gocql.UUID{}, "", "", err
	}
	return id, email, role, nil
}

func (r *Repository) UpdateThirdPartyEmail(ctx context.Context, email string, thirdPartyId string, providerName string) error {
	err := r.session.Query(
		`UPDATE member_by_thrid_party SET email = ? WHERE third_party_id = ? AND provider_name = ?`,
		email, thirdPartyId, providerName).
		ExecContext(ctx)
	if err != nil {
		slog.Error("fail to update third party email",
			"err", err,
			"email", email,
			"thirdPartyId", thirdPartyId,
			"providerName", providerName)
		return err
	}
	return nil
}
