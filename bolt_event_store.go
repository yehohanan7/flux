package cqrs

import (
	"bytes"
	"log"

	"strconv"

	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

var BUCKET_NAME = []byte("Events")

type BoltEventStore struct {
	db *bolt.DB
}

func buildKey(aggregateId string, event Event) string {
	return strings.Join([]string{aggregateId, strconv.Itoa(event.AggregateVersion)}, "_")
}

func (store *BoltEventStore) GetEvents(aggregateId string) []Event {
	events := make([]Event, 0)

	store.db.View(func(tx *bolt.Tx) error {

		c := tx.Bucket([]byte(BUCKET_NAME)).Cursor()
		prefix := []byte(aggregateId + "_")

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			event := new(Event)
			event.Deserialize(v)
			events = append(events, *event)
		}
		return nil
	})

	return events
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []Event) error {

	store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BUCKET_NAME)

		for _, event := range events {
			b.Put([]byte(buildKey(aggregateId, event)), event.Serialize())
		}
		return nil
	})

	return nil
}

func NewBoltEventStore(path string) EventStore {

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		glog.Errorf("error opening bolt store %v", path)
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(BUCKET_NAME)
		glog.Errorf("error while creating bucket %s", BUCKET_NAME)
		return err
	})

	return &BoltEventStore{db}
}
