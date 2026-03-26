package repository

import (
	"backend/auth/internal/constant"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) FindRefreshTokenJTIById(id gocql.UUID) (jti gocql.UUID, err error) {
	err = r.session.Query(
		"SELECT refresh_token_jti from member_by_id WHERE id = ?",
		id,
	).Scan(&jti)
	if err != nil {
		slog.Info("fail to get refresh token jti",
			"err", err,
		)
		return gocql.UUID{}, err
	}
	return jti, nil
}

func (r *Repository) SaveRefreshTokenJTIById(id, jti gocql.UUID) error {
	err := r.session.Query(
		"UPDATE member_by_id USING TTL ? SET refresh_token_jti = ? WHERE id = ?",
		constant.RefreshTokenTTL, jti, id,
	).Exec()
	if err != nil {
		slog.Error("fail to save refresh token jti",
			"err", err,
		)
		return err
	}
	return nil
}
