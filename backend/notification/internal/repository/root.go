package repository

import (
	"backend/common"
	"log"
	"os"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
)

type Repository struct {
	session *gocql.Session
}

func NewRepository() *Repository {
	k8ssandraHost := os.Getenv("K8SSANDRA_HOST")
	cluster := gocql.NewCluster(k8ssandraHost)
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
	tlSConfig.ServerName = k8ssandraHost
	cluster.SslOpts = &gocql.SslOptions{
		Config:                 tlSConfig,
		EnableHostVerification: true,
	}

	session, err := gocql.NewSession(*cluster)
	if err != nil {
		panic(err)
	}

	err = session.Query(`CREATE TABLE IF NOT EXISTS notification_info_by_id (
    id uuid,
    os text,
    device_push_token text,
    PRIMARY KEY (id, device_push_token));`).Exec()
	if err != nil {
		panic(err)
	}
	err = session.Query(`CREATE TABLE IF NOT EXISTS id_by_device_push_token (
    device_push_token text,
    id uuid,
    PRIMARY KEY (device_push_token));`).Exec()
	if err != nil {
		panic(err)
	}

	log.Print("success to connect cassandra")
	r := &Repository{
		session: session,
	}

	return r
}
