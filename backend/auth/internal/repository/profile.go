package repository

import (
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) SaveProfileId(id gocql.UUID) error {
	err := r.session.Query("INSERT INTO profile_by_id (id) VALUES (?)", id).Exec()
	if err != nil {
		slog.Error("fail to insert profile id",
			"err", err,
			"id", id.String())
		return err
	}
	return nil
}
