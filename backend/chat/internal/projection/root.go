package projection

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type FindMessagesByToIdAndId struct {
	Id          gocql.UUID
	RoomId      []byte
	FromId      gocql.UUID
	ContentType string
	Contents    []string
}
