package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
)

const (
	EVENTS_BUCKET         = "EVENTS"
	EVENT_METADATA_BUCKET = "EVENT_METADATA"
)

//InMemory implementation of the event store
type BoltEventStore struct {
	db *bolt.DB
}

func (store *BoltEventStore) GetEvent(id string) Event {
	var event = new(Event)
	store.db.View(func(tx *bolt.Tx) error {
		eventsBucket := tx.Bucket([]byte(EVENTS_BUCKET))
		return fetch(eventsBucket, []byte(id), event)
	})
	return *event
}

func (store *BoltEventStore) GetEvents(aggregateId string) []Event {
	return nil
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []Event) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		eventsBucket := tx.Bucket([]byte(EVENTS_BUCKET))
		for _, event := range events {
			if err := save(eventsBucket, []byte(event.Id), event); err != nil {
				return err
			}
		}
		return nil
	})
}

func (store *BoltEventStore) GetEventMetaDataFrom(offset, count int) []EventMetaData {
	return nil
}

func NewBoltStore(path string) *BoltEventStore {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		glog.Fatal("Error while opening bolt db", err)
	}
	db.Update(func(tx *bolt.Tx) error {
		createBucket(tx, EVENTS_BUCKET)
		createBucket(tx, EVENT_METADATA_BUCKET)
		return nil
	})
	return &BoltEventStore{db}
}
