package repository

import (
	"log"
	"os"
	"time"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/apache/cassandra-gocql-driver/v2/lz4"
	gocqlastra "github.com/datastax/gocql-astra/v2"

	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
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
	log.Print("success to connect cassandra")
	r := &Repository{
		session: session,
	}

	return r
}
