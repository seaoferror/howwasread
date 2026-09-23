package projection

import (
	"time"
)

type Detail struct {
	Novel, ShortStory, Poem, Play, Film, WrittenBy, Rule         string
	Capacity, LengthMinutes                                      int
	UpdatedAt, Time                                              time.Time
	IsModerator, IsRegistrant, IsBanned, IsNotificationScheduled bool
}
