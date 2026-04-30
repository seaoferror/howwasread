package repository

import (
	"log"
	"os"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
	session *gocql.Session
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
	err = session.Query(`INSERT INTO chat_room_by_id (id, name, room_type) 
	VALUES (
		019e0e84-f358-71a2-8b3c-d4e5f6012345, 
	    'yuan',
		'personal'
	);`).Exec()
	if err != nil {
		log.Panicf("fail to insert dummy chat_room_by_id data: %v", err)
	}

	log.Print("success to connect cassandra")

	r := Repository{
		session: session,
	}

	return &r
}
