package repository

import (
	"log"
	"os"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/rueidis"
)

type Repository struct {
	session *gocql.Session
	client  rueidis.Client
}

func NewRepository() *Repository {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_URL"))
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
	err = session.Query(`CREATE TABLE IF NOT EXISTS message_by_to_id (
    to_id_type text,
    to_id uuid,
    read boolean,
    id uuid,
    from_id uuid,
    content_type text,
    content text,
    PRIMARY KEY ((to_id_type, to_id), read, id)
    ) WITH CLUSTERING ORDER BY (read ASC, id DESC);`).Exec()
	if err != nil {
		log.Panicf("fail to create table payload: %v", err)
	}

	log.Print("success to connect cassandra")

	clientOption := rueidis.ClientOption{
		InitAddress: []string{os.Getenv("REDIS_URL")},
	}
	client, err := rueidis.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}

	r := &Repository{
		session: session,
		client:  client,
	}

	return r
}
