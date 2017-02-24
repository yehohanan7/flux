package flux

import (
	gorillaMux "github.com/gorilla/mux"
	. "github.com/yehohanan7/flux/cqrs"
	"github.com/yehohanan7/flux/memory"
	"github.com/yehohanan7/flux/mux"
)

func NewEventStore() EventStore {
	return memory.NewEventStore()
}

func StartMuxEventFeed(router *gorillaMux.Router, store EventStore, eventsPerPage ...int) {
	mux.EventFeed(router, store)
}
