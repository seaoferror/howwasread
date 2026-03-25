package document

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Member is club independent so reference
// since we use uuid in token claim subject we use bson.Binary instead of bson.ObjectID
type Member struct {
	Id       bson.Binary `bson:"_id"`
	Name     string      `bson:"name"`
	ServerIP string      `bson:"server_ip"`

	ModeratorConversationIds  []bson.ObjectID `bson:"m_c_ids"`
	RegistrantConversationIds []bson.ObjectID `bson:"r_c_ids"`
}

type Conversation struct {
	Id         bson.ObjectID `bson:"_id"`
	Novel      string        `bson:"novel,omitempty"`
	ShortStory string        `bson:"short_story,omitempty"`
	Poem       string        `bson:"poem,omitempty"`
	Play       string        `bson:"play,omitempty"`
	Film       string        `bson:"film,omitempty"`
	By         string        `bson:"by,omitempty"`
	Rule       string        `bson:"rule,omitempty"`
	Capacity   int           `bson:"capacity"`
	When       time.Time     `bson:"when"`
	Length     time.Duration `bson:"length"`
	Expired    bool          `bson:"expired"`

	ModeratorIds   []bson.Binary `bson:"m_ids"`
	RegistrantIds  []bson.Binary `bson:"r_ids"`
	ParticipantIds []bson.Binary `bson:"p_ids"`
	BanIds         []bson.Binary `bson:"b_ids"`
}

//type Org struct {
//	Id          bson.ObjectID `bson:"_id"`
//	Name        string        `bson:"name"`
//	Description string        `bson:"description"`
//
//	ModIds        []bson.Binary `bson:"mod_ids"`
//	SubscriberIds []bson.Binary `bson:"subscriber_ids"`
//	//this is for blocking participation at the further meetings,
//	//listening(subscribing) can't be blocked
//	BlockedIds []bson.Binary `bson:"blocked_ids"`
//
//	MeetingIds []bson.ObjectID `bson:"meeting_ids"`
//}

// QNA is meeting subset so it seemed good to embedding,
// but embedding make query writing difficult, so we will ref it
//type QNA struct {
//	Id       bson.ObjectID `bson:"_id"`
//	Question string        `bson:"question"`
//
//	MeetingId bson.ObjectID   `bson:"meeting_id"`
//	AnswerIds []bson.ObjectID `bson:"answer_ids"`
//}
//
//type Answer struct {
//	Id      bson.ObjectID `bson:"_id"`
//	Content string        `bson:"content"`
//}
