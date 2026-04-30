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
	err = session.Query(`CREATE TABLE IF NOT EXISTS profile_by_id (
    id uuid,
    name text,
    PRIMARY KEY (id)
    );`).Exec()
	if err != nil {
		log.Panicf("fail to create table profile_by_id: %v", err)
	}
	err = session.Query(`INSERT INTO profile_by_id (id, name) 
	VALUES (
		019e0e84-f358-71a2-8b3c-d4e5f6012345, 
		'yuan'
	);`).Exec()
	if err != nil {
		log.Panicf("fail to insert dummy profile_by_id data: %v", err)
	}
	err = session.Query(`INSERT INTO profile_by_id (id, name) 
	VALUES (
		019e0e84-f358-7d8e-9fa0-b1c2d3e4f506, 
		'yen'
	);`).Exec()
	if err != nil {
		log.Panicf("fail to insert dummy profile_by_id data: %v", err)
	}

	log.Print("success to connect cassandra")

	clientOption := rueidis.ClientOption{
		InitAddress: []string{os.Getenv("REDIS_URL")},
	}
	client, err := rueidis.NewClient(clientOption)
	if err != nil {
		log.Panicf("Fail to connect to redis: %v", err)
	}
	log.Print("success to connect redis")

	r := &Repository{
		session: session,
		client:  client,
	}

	return r
}
