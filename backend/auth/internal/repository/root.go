package repository

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	gocqlastra "github.com/datastax/gocql-astra/v2"

	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
	session *gocql.Session
}

func NewRepository() *Repository {
	cluster, err := gocqlastra.NewClusterFromBundle("cert/astradb/astradb-secure-connect.zip",
		"token", os.Getenv("ASTRADB_TOKEN"), 10*time.Second)
	cluster.Keyspace = os.Getenv("PROFILE")
	cluster.Timeout = 1 * time.Minute
	cluster.Consistency = gocql.Quorum
	cluster.Compressor = &lz4.LZ4Compressor{}
	//cluster.PageSize = 1000
	//cluster.NextPagePrefetch = 0.25
	//cluster.Tracer =

	session, err := gocql.NewSession(*cluster)

	queries := []string{
		`CREATE TABLE IF NOT EXISTS member_by_email (
			email text PRIMARY KEY, id uuid, email_verified boolean, phone_number_verified boolean,
			phone_number text, password text, role text, apple_sign_in_user text
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_id (
			id uuid PRIMARY KEY, email text, email_verified boolean, phone_number_verified boolean,
			phone_number text, password text, role text, apple_sign_in_user text, refresh_token_jti uuid
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_phone_number (
			phone_number text PRIMARY KEY, id uuid, email text, phone_number_verified boolean, role text
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_apple_sign_in_user (
			apple_sign_in_user text PRIMARY KEY, id uuid, email text, email_verified boolean, role text
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_verification_id (
			verification_id uuid PRIMARY KEY, email text, phone_number text, otp text
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_session_id (
			session_id uuid PRIMARY KEY, email text
		);`,
		`CREATE TABLE IF NOT EXISTS profile_by_id (
			id uuid PRIMARY KEY
		);`,
	}

	for _, q := range queries {
		err = session.Query(q).Exec()
		if err != nil {
			slog.Error("failed to execute migration query", "err", err, "query", q)
			panic(err)
		}
	}

	if err != nil {
		log.Panicf("fail to create session from cassandra cluster: %v", err)
	}
	log.Print("success to connect cassandra")
	r := &Repository{
		session: session,
	}

	return r
}
