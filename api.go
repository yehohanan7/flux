package flux

import (
	gorillaMux "github.com/gorilla/mux"
	"github.com/yehohanan7/flux/consumer"
	. "github.com/yehohanan7/flux/cqrs"
	"github.com/yehohanan7/flux/memory"
	"github.com/yehohanan7/flux/mux"
)

//Create a new in memory store
func NewEventStore() EventStore {
	return memory.NewEventStore()
}

//Start feeding events over the mux router
func StartMuxEventFeed(router *gorillaMux.Router, store EventStore, eventsPerPage ...int) {
	mux.EventFeed(router, store)
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
