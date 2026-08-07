package repository

import (
	"backend/common/tlsconfig"
	"log"
	"os"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	gocqlastra "github.com/datastax/gocql-astra/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/valkey-io/valkey-go"
)

type Repository struct {
	client  valkey.Client
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
	if err != nil {
		log.Panicf("fail to create session from cassandra cluster: %v", err)
	}
	err = session.Query(`CREATE TABLE IF NOT EXISTS ids_by_filename (
    filename uuid,
    ids set<uuid>,
    PRIMARY KEY (filename)
    );`).Exec()
	if err != nil {
		log.Panicf("fail to create table ids_by_filename: %v", err)
	}
	err = session.Query(`CREATE TABLE IF NOT EXISTS block (
    blocker_id uuid,
    blocked_id uuid,
    PRIMARY KEY ((blocker_id), blocked_id)
	);`).Exec()
	if err != nil {
		log.Panicf("fail to create table block: %v", err)
	}

	log.Print("success to connect cassandra")

	clientOption := valkey.ClientOption{
		InitAddress: []string{os.Getenv("VALKEY_ADDRESS")},
	}

	tlsConfig, err := tlsconfig.Create("", "", os.Getenv("VALKEY_CA_CERT_PATH"))
	if err != nil {
		log.Panicf("fail to create tls config for valkey")
	}
	tlsConfig.ServerName = os.Getenv("VALKEY_HOST")
	clientOption.TLSConfig = tlsConfig

	clientOption.Username = os.Getenv("VALKEY_USERNAME")
	clientOption.Password = os.Getenv("VALKEY_PASSWORD")

	client, err := valkey.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := Repository{
		session: session,
		client:  client,
	}

	return &r
}
