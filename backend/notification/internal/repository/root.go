package repository

import (
	"log"
	"os"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	gocqlastra "github.com/datastax/gocql-astra/v2"
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
	if err != nil {
		log.Panicf("fail to create session from cassandra cluster: %v", err)
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
	err = session.Query(`CREATE TABLE conversation_notification_by_scheduled_time (
    scheduled_time timestamp,
    conversation_type text,
    conversation_id uuid,
    member_id uuid,
    PRIMARY KEY ((scheduled_time), conversation_type, conversation_id, member_id));`).Exec()
	if err != nil {
		panic(err)
	}

	log.Print("success to connect cassandra")
	r := &Repository{
		session: session,
	}

	return r
}
