package repository

import (
	"backend/common"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
	session *gocql.Session
}

func NewRepository() *Repository {
	cluster := gocql.NewCluster(os.Getenv("K8SSANDRA_HOST"))
	cluster.Port = 9042
	cluster.Keyspace = os.Getenv("PROFILE")
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("K8SSANDRA_USERNAME"),
		Password: os.Getenv("K8SSANDRA_PASSWORD"),
	}
	cluster.Timeout = 1 * time.Minute
	cluster.Consistency = gocql.LocalOne
	cluster.Compressor = &lz4.LZ4Compressor{}
	cluster.PageSize = 1000
	cluster.NextPagePrefetch = 0.25
	tlSConfig, err := common.CreateTlSConfig("", "", os.Getenv("K8SSANDRA_CA_CERT_PATH"))
	if err != nil {
		panic(err)
	}
	cluster.SslOpts = &gocql.SslOptions{
		Config:                 tlSConfig,
		EnableHostVerification: false,
	}

	session, err := gocql.NewSession(*cluster)
	if err != nil {
		panic(err)
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS member_by_email (
			email text PRIMARY KEY, id uuid, email_verified boolean, phone_number_verified boolean,
			phone_number text, password text, role text
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_id (
			id uuid PRIMARY KEY, email text, email_verified boolean, phone_number_verified boolean,
			phone_number text, role text, refresh_token_jtis set<uuid>
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_phone_number (
			phone_number text PRIMARY KEY, id uuid, email text, phone_number_verified boolean, role text
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_verification_id (
			verification_id uuid PRIMARY KEY, email text, phone_number text, otp text
		);`,
		`CREATE TABLE IF NOT EXISTS member_by_session_id (
			session_id uuid PRIMARY KEY, email text
		);`,
		`CREATE TABLE IF NOT EXISTS nonce (
    		nonce text PRIMARY KEY
    	);`,
	}

	for _, q := range queries {
		err = session.Query(q).Exec()
		if err != nil {
			slog.Error("failed to execute migration query", "err", err, "query", q)
			panic(err)
		}
	}
	log.Print("success to connect cassandra")
	r := &Repository{
		session: session,
	}

	return r
}
