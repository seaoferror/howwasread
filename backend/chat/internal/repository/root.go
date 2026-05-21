package repository

import (
	"log"
	"os"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	gocqlastra "github.com/datastax/gocql-astra/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/rueidis"
)

type Repository struct {
	session *gocql.Session
	client  rueidis.Client
}

func NewRepository() *Repository {
	cluster, err := gocqlastra.NewClusterFromBundle("/astradb/astradb-secure-connect.zip",
		"token", os.Getenv("ASTRA_DB_TOKEN"), 10*time.Second)
	cluster.Keyspace = "default"
	cluster.Timeout = 1 * time.Minute
	cluster.Consistency = gocql.Quorum
	cluster.Compressor = &lz4.LZ4Compressor{}
	//cluster.Authenticator = gocql.PasswordAuthenticator{
	//	Username: "",
	//	Password: "",
	//}
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
	err = session.Query(`INSERT INTO chat_room_by_id (id, name, room_type) 
	VALUES (
		019e0e84-f358-71a2-8b3c-d4e5f6012345, 
	    'yuan',
		'personal'
	);`).Exec()
	if err != nil {
		log.Panicf("fail to insert dummy chat_room_by_id data: %v", err)
	}
	err = session.Query(`INSERT INTO chat_room_by_id (id, name, room_type) 
	VALUES (
		019e0e84-f358-7d8e-9fa0-b1c2d3e4f506, 
	    'yen',
		'personal'
	);`).Exec()
	if err != nil {
		log.Panicf("fail to insert dummy chat_room_by_id data: %v", err)
	}
	err = session.Query(`INSERT INTO message_by_to_id (id, to_id, from_id, room_id, content_type, contents) 
	VALUES (
		018f4b1a-e6b0-7000-811c-cdef01234567, 
		019e0e84-f358-7d8e-9fa0-b1c2d3e4f506,
		019e0e84-f358-71a2-8b3c-d4e5f6012345,
		019e0e84-f358-71a2-8b3c-d4e5f6012345,
	    'text',
		{'Hello Yen, this is Yuan!'}
	);`).Exec()
	if err != nil {
		log.Panicf("fail to insert dummy message (yuan -> yen): %v", err)
	}

	// Message from 'yen' to 'yuan'
	err = session.Query(`INSERT INTO message_by_to_id (id, to_id, from_id, room_id, content_type, contents) 
	VALUES (
		018f4b1a-e6b0-7000-811c-cdef01234567, 
		019e0e84-f358-71a2-8b3c-d4e5f6012345, 
		019e0e84-f358-7d8e-9fa0-b1c2d3e4f506, 
		019e0e84-f358-7d8e-9fa0-b1c2d3e4f506,
	    'text',
		{'Hi Yuan, great to hear from you!'}
	);`).Exec()
	if err != nil {
		log.Panicf("fail to insert dummy message (yen -> yuan): %v", err)
	}
	log.Print("success to connect cassandra")

	clientOption := rueidis.ClientOption{
		InitAddress: []string{os.Getenv("REDIS_ADDRESS")},
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
