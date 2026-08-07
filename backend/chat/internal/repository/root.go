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
	session *gocql.Session
	client  valkey.Client
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

	err = session.Query(`CREATE TABLE IF NOT EXISTS message_by_to_id (
    id uuid,
    to_id uuid,
    from_id uuid,
    room_id uuid,
    content_type text,
    contents set<text>,
    PRIMARY KEY ((to_id), id)
    ) WITH CLUSTERING ORDER BY (id DESC);`).Exec()
	if err != nil {
		log.Panicf("fail to create table payload: %v", err)
	}
	err = session.Query(`CREATE TABLE IF NOT EXISTS profile_by_id (
    id uuid,
    name text,
    reporter_ids set<uuid>,
    blocked_conversations set<uuid>,
    PRIMARY KEY (id)
    );`).Exec()
	if err != nil {
		log.Panicf("fail to create table profile_by_id: %v", err)
	}
	err = session.Query(`CREATE TABLE IF NOT EXISTS chat_room_by_id (
    id uuid,
    name text,
    room_type text,
    participant_ids set<uuid>,
    PRIMARY KEY (id)
    );`).Exec()
	if err != nil {
		log.Panicf("fail to create table chat_room_by_id: %v", err)
	}
	err = session.Query(`CREATE TABLE IF NOT EXISTS banned_phone_number (
    phone_number text,
    PRIMARY KEY (phone_number)
    );`).Exec()
	if err != nil {
		log.Panicf("fail to create table chat_room_by_id: %v", err)
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
		log.Panicf("Fail to connect to valkey: %v", err)
	}
	log.Print("success to connect valkey")

	r := &Repository{
		session: session,
		client:  client,
	}

	return r
}
