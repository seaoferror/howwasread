package repository

import (
	"crypto/tls"
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
	cluster, err := gocqlastra.NewClusterFromBundle("/astradb/astradb-secure-connect.zip",
		"token", os.Getenv("ASTRA_DB_TOKEN"), 10*time.Second)
	cluster.Keyspace = "default"
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

	log.Print("success to connect cassandra")

	clientOption := valkey.ClientOption{
		Username:    os.Getenv("VALKEY_USERNAME"),
		Password:    os.Getenv("VALKEY_PASSWORD"),
		InitAddress: []string{os.Getenv("VALKEY_ADDRESS")},
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
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
