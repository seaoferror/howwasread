package document

import "go.mongodb.org/mongo-driver/v2/bson"

type Member struct {
	Id       bson.Binary `bson:"_id"`
	Name     string      `bson:"name"`
	ServerIP string      `bson:"server_ip"`

	ModeratorConversationIds  []bson.ObjectID `bson:"m_c_ids"`
	RegistrantConversationIds []bson.ObjectID `bson:"r_c_ids"`
}
