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

//Start feeding events over the mux router
func FeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return feed.FeedHandler(store)
}

//Create new consumer
func NewEventConsumer(url string, handlerClass interface{}, store ...OffsetStore) EventConsumer {
	if len(store) == 0 {
		return consumer.NewEventConsumer(url, handlerClass, memory.NewOffsetStore())
	}
	return consumer.NewEventConsumer(url, handlerClass, store[0])
}

//Create new in memory offset store
func NewOffsetStore() OffsetStore {
	return memory.NewOffsetStore()
}
