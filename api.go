package flux

import (
	"net/http"
	"time"

	"github.com/yehohanan7/flux/consumer"
	. "github.com/yehohanan7/flux/cqrs"
	"github.com/yehohanan7/flux/feed"
	"github.com/yehohanan7/flux/memory"
)

//Create a new in memory store
func NewEventStore() EventStore {
	return memory.NewEventStore()
}

//Create new in memory offset store
func NewOffsetStore() OffsetStore {
	return memory.NewOffsetStore()
}

//Get event feed
func GetFeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return feed.GetFeedHandler(store)
}

//Get event by event id
func GetEventHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return feed.GetEventHandler(store)
}

//Create new consumer
func NewEventConsumer(url string, events []interface{}, store OffsetStore) EventConsumer {
	return consumer.New(url, events, store, 5*time.Second)
}
