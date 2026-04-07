package data

import gocql "github.com/apache/cassandra-gocql-driver/v2"

type FindDevicePushToken struct {
	Id              gocql.UUID
	DeviceId        gocql.UUID
	OS              string
	DevicePushToken string
}
