package boltdb

import (
	"bytes"
	"log"

	"strconv"

	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/yehohanan7/cqrs"
)

var BUCKET_NAME = []byte("Events")

type BoltEventStore struct {
	db *bolt.DB
}

func buildKey(aggregateId string, event cqrs.Event) string {
	return strings.Join([]string{aggregateId, strconv.Itoa(event.AggregateVersion)}, "_")
}

func (store *BoltEventStore) GetEvents(aggregateId string) []cqrs.Event {
	events := make([]cqrs.Event, 0)

	store.db.View(func(tx *bolt.Tx) error {

		c := tx.Bucket([]byte(BUCKET_NAME)).Cursor()
		prefix := []byte(aggregateId + "_")

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			event := new(cqrs.Event)
			event.Deserialize(v)
			events = append(events, *event)
		}
		return nil

	})

	return events
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []cqrs.Event) error {

	store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BUCKET_NAME)

		for _, event := range events {
			b.Put([]byte(buildKey(aggregateId, event)), event.Serialize())
		}
		return nil
	})

	return nil
}

func (store *BoltEventStore) GetAllEvents() []cqrs.Event {
	events := make([]cqrs.Event, 0)

	store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_NAME))

		b.ForEach(func(k, v []byte) error {
			event := new(cqrs.Event)
			event.Deserialize(v)
			events = append(events, *event)
			return nil
		})
		return nil
	})

	return events
}

func (store *BoltEventStore) GetEvent(id string) cqrs.Event {
	var result cqrs.Event

	store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_NAME))

		b.ForEach(func(k, v []byte) error {
			event := new(cqrs.Event)
			event.Deserialize(v)
			if event.Id == id {
				result = *event
				return nil
			}
			return nil
		})
		return nil
	})

	return result
}

func NewEventStore(path string) cqrs.EventStore {

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		glog.Errorf("error opening bolt store %v", path)
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(BUCKET_NAME)
		if err != nil {
			glog.Error("error while creating bucket", BUCKET_NAME)
		}
		return err
	})

	return &BoltEventStore{db}
}
