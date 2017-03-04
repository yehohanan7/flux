package feed

import (
	"fmt"
	"reflect"

	. "github.com/yehohanan7/flux/cqrs"
)

type JsonEventFeed struct {
	Description string       `json:"description"`
	Events      []EventEntry `json:"events"`
}

type EventEntry struct {
	EventId          string `json:"event_id"`
	Url              string `json:"url"`
	AggregateName    string `json:"aggregate_name"`
	AggregateVersion int    `json:"aggregate_version"`
	EventType        string `json:"event_type"`
	Created          string `json:"created"`
}

func ToEventEntry(url string, event Event) EventEntry {
	return EventEntry{
		event.Id,
		fmt.Sprintf("%s/%s", url, event.Id),
		event.AggregateName,
		event.AggregateVersion,
		reflect.TypeOf(event.Payload).String(),
		event.OccuredAt,
	}
}
