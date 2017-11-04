package flux

import (
	"net/http"
	"time"

	"github.com/yehohanan7/flux/mongodb"

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

//Create a new mongo store
func NewMongoStore(options *mongodb.MongoEventStoreOptions) EventStore {
	return mongodb.NewEventStore(options)
}

//Create new in memory offset store
func NewMemoryOffsetStore() OffsetStore {
	return memory.NewOffsetStore()
}

//Offset store backed by boltdb
func NewBoltOffsetStore(path string) OffsetStore {
	return boltdb.NewOffsetStore(path)
}

//Offset store backed by mongodb
func NewMongoOffsetStore(options *mongodb.MongoOffsetStoreOptions) OffsetStore {
	return mongodb.NewOffsetStore(options)
}

//Start feeding events over the mux router
func FeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return feed.FeedHandler(store)
}

//Create new consumer
func NewEventConsumer(url string, interval time.Duration, events []interface{}, store OffsetStore) EventConsumer {
	return consumer.New(url, events, store, interval)
}
