package flux

import (
	"net/http"

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

//Start feeding events over the mux router
func FeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return feed.FeedHandler(store)
}

//Create new consumer
func NewEventConsumer(url string, events []interface{}, store OffsetStore) EventConsumer {
	return consumer.New(url, events, store)
}
