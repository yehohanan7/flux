package cqrs

import (
	"fmt"
	"log"

	"strconv"

	"github.com/boltdb/bolt"
)

var BUCKET_NAME = []byte("Events")

type BoltEventStore struct {
	db *bolt.DB
}

func (store *BoltEventStore) GetEvents(aggregateId string) []Event {
	events := make([]Event, 0)

	store.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(BUCKET_NAME).Bucket([]byte(aggregateId)); b != nil {
			b.ForEach(func(k, v []byte) error {
				event := new(Event)
				event.Deserialize(v)
				events = append(events, *event)
				return nil
			})
		}
		return nil
	})

	return events
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []Event) error {

	store.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket(BUCKET_NAME).CreateBucketIfNotExists([]byte(aggregateId))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for _, event := range events {
			key := []byte(strconv.Itoa(event.AggregateVersion))
			b.Put(key, event.Serialize())
		}
		return nil
	})

	return nil
}

func NewBoltEventStore(path string) EventStore {

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(BUCKET_NAME)
		//TODO: add warnings
		return err
	})

	return &BoltEventStore{db}
}
