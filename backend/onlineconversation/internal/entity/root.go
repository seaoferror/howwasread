package entity

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	Id              uuid.UUID
	Novel           string
	ShortStory      string
	Poem            string
	Play            string
	Film            string
	WrittenBy       string
	Rule            string
	Capacity        int
	Time            time.Time
	Length          time.Duration
	ModeratorIds    []uuid.UUID
	RegistrantIds   []uuid.UUID
	BanIds          []uuid.UUID
	ReporterIds     []uuid.UUID
	NotificationIds []uuid.UUID
}

//types Org struct {
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
//types QNA struct {
//	Id       bson.ObjectID `bson:"_id"`
//	Question string        `bson:"question"`
//
//	MeetingId bson.ObjectID   `bson:"meeting_id"`
//	AnswerIds []bson.ObjectID `bson:"answer_ids"`
//}
//
//types Answer struct {
//	Id      bson.ObjectID `bson:"_id"`
//	Content string        `bson:"content"`
//}
