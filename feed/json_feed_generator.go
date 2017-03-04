package feed

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
)

type JsonFeedGenerator struct {
}

func (_ JsonFeedGenerator) Generate(url, description string, events []Event) []byte {

	entries := make([]EventEntry, 0)

	for _, event := range events {
		entries = append(entries, EventEntry{
			event.Id,
			fmt.Sprintf("%s/%s", url, event.Id),
			event.AggregateName,
			event.AggregateVersion,
			reflect.TypeOf(event.Payload).String(),
			event.OccuredAt,
		})
	}

	jsonFeed := JsonEventFeed{description, entries}

	b, err := json.Marshal(jsonFeed)
	if err != nil {
		glog.Warning("error while creating json for ", events, err)
		return []byte{}
	}

	return b
}
