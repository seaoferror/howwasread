package projection

import gocql "github.com/apache/cassandra-gocql-driver/v2"

type FindPushTokensById struct {
	Id              gocql.UUID
	OS              string
	DevicePushToken string
}
