package repository

import (
	"crypto/tls"
	"crypto/x509"
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
	if os.Getenv("PROFILE") == "production" {
		caCertPEM, err1 := os.ReadFile("/cert/valkey/ca.crt")
		if err1 != nil {
			log.Fatalf("Failed to read CA certificate: %v", err1)
		}
		rootCAs := x509.NewCertPool()
		ok := rootCAs.AppendCertsFromPEM(caCertPEM)
		if !ok {
			log.Fatalf("Failed to parse root certificate")
		}
		clientOption.Username = os.Getenv("VALKEY_USERNAME")
		clientOption.Password = os.Getenv("VALKEY_PASSWORD")
		clientOption.TLSConfig = &tls.Config{
			RootCAs:            rootCAs,
			InsecureSkipVerify: false,
			ServerName:         os.Getenv("VALKEY_HOST"),
		}
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
