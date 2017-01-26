package cqrs

import (
	"encoding/json"
	"fmt"
	"reflect"

	"strconv"

	"github.com/golang/glog"
)

type JsonEventFeed struct {
	Description string       `json:"description"`
	Event       []EventEntry `json:"events"`
}

type EventEntry struct {
	EventId          string `json:"event_id"`
	Url              string `json:"url"`
	AggregateName    string `json:"aggregate_name"`
	AggregateVersion string `json:"aggregate_version"`
	EventType        string `json:"event_type"`
	Created          string `json:"created"`
}

type JsonFeedGenerator struct {
}

func (_ JsonFeedGenerator) Generate(url, description string, events []Event) string {

	entries := make([]EventEntry, 0)

	for _, event := range events {
		entries = append(entries, EventEntry{
			event.Id,
			fmt.Sprintf("%s/%s", url, event.Id),
			event.AggregateName,
			strconv.Itoa(event.AggregateVersion),
			reflect.TypeOf(event.Payload).String(),
			event.OccuredAt,
		})
	}

	jsonFeed := JsonEventFeed{description, entries}

	b, err := json.Marshal(jsonFeed)
	if err != nil {
		glog.Warning("error while creating json for ", events, err)
		return ""
	}

	return string(b)
}
