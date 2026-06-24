package repository

import (
	"backend/auth/internal/constant"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) FindRefreshTokenJTIsById(id gocql.UUID) (jtis []gocql.UUID, err error) {
	err = r.session.Query(
		"SELECT refresh_token_jtis FROM member_by_id WHERE id = ?",
		id,
	).Scan(&jtis)
	if err != nil {
		slog.Info("fail to get refresh token jti",
			"err", err,
		)
		return nil, err
	}
	return jtis, nil
}

func (r *Repository) SaveRefreshTokenJTIById(id, jti gocql.UUID) error {
	err := r.session.Query(
		"UPDATE member_by_id USING TTL ? SET refresh_token_jtis += ? WHERE id = ?",
		constant.RefreshTokenTTL, []gocql.UUID{jti}, id,
	).Exec()
	if err != nil {
		slog.Error("fail to save refresh token jti",
			"err", err,
		)
		return err
	}
	return nil
}

func (r *Repository) DeleteAccount(id gocql.UUID) error {
	err := r.session.Batch(gocql.CounterBatch).
		Query("DELETE FROM member_by_id WHERE id = ?", id).
		Query("DELETE FROM profile_by_id WHERE id = ?", id).
		Exec()
	if err != nil {
		slog.Error("fail to delete account",
			"err", err,
			"id", id.String(),
		)
		return err
	}
	return nil
}
