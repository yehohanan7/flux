package cqrs

import (
	"fmt"
	"log"

	"encoding/json"

	"strconv"

	"github.com/boltdb/bolt"

	. "github.com/yehohanan7/cqrs/event"
)

var BUCKET_NAME = []byte("Events")

type BoltEventStore struct {
	db *bolt.DB
}

func (store *BoltEventStore) GetEvents(aggregateId string) []Event {
	events, event := make([]Event, 0), new(Event)

	store.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(BUCKET_NAME).Bucket([]byte(aggregateId)); b != nil {
			b.ForEach(func(k, v []byte) error {
				if err := json.Unmarshal(v, event); err != nil {
					return fmt.Errorf("error while unmarshalling %s", err)
				}
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
			value, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("error marshaling event %s", err)
			}
			b.Put(key, value)
		}
		return nil
	})

	return nil
}

func NewBoltEventStore(path string) EventStore {
	var (
		db  *bolt.DB
		err error
	)

	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(BUCKET_NAME)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &BoltEventStore{db}
}
