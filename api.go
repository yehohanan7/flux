package flux

import (
	"net/http"
	"time"

	"github.com/yehohanan7/flux/boltdb"
	"github.com/yehohanan7/flux/consumer"
	. "github.com/yehohanan7/flux/cqrs"
	"github.com/yehohanan7/flux/feed"
	"github.com/yehohanan7/flux/memory"
)

//Create a new memory store
func NewMemoryStore() EventStore {
	return memory.NewEventStore()
}

//Create a new bolt store
func NewBoltStore(path string) EventStore {
	return boltdb.NewBoltStore(path)
}

//Create new in memory offset store
func NewMemoryOffsetStore() OffsetStore {
	return memory.NewOffsetStore()
}

//Offset store backed by boltdb
func NewBoltOffsetStore(path string) OffsetStore {
	return boltdb.NewOffsetStore(path)
}

//Start feeding events over the mux router
func FeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return feed.FeedHandler(store)
}

//Create new consumer
func NewEventConsumer(url string, events []interface{}, store OffsetStore) EventConsumer {
	return consumer.New(url, events, store, 5*time.Second)
}
