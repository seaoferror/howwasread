package data

import gocql "github.com/apache/cassandra-gocql-driver/v2"

type FindPushTokensById struct {
	Id              gocql.UUID
	DeviceId        gocql.UUID
	OS              string
	DevicePushToken string
}
