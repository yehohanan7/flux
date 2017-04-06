package boltdb

import (
	"bytes"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
)

const (
	EVENTS_BUCKET         = "EVENTS"
	EVENT_METADATA_BUCKET = "EVENT_METADATA"
	AGGREGATES_BUCKET     = "AGGREGATES_BUCKET"
)

//InMemory implementation of the event store
type BoltEventStore struct {
	db *bolt.DB
}

func (store *BoltEventStore) GetEvent(id string) Event {
	var event = new(Event)
	store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EVENTS_BUCKET))

		if v := b.Get([]byte(id)); v != nil {
			if err := event.Deserialize(v); err != nil {
				glog.Errorf("could not deserialize event with id %s", id)
				return err
			}
		}
		return nil
	})
	return *event
}

func (store *BoltEventStore) GetEvents(aggregateId string) []Event {
	events := make([]Event, 0)
	store.db.View(func(tx *bolt.Tx) error {

		c := tx.Bucket([]byte(AGGREGATES_BUCKET)).Cursor()

		prefix := []byte(aggregateId)

		for k, eventId := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, eventId = c.Next() {
			events = append(events, store.GetEvent(string(eventId)))
		}

		return nil
	})

	return events
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []Event) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		eventsBucket := tx.Bucket([]byte(EVENTS_BUCKET))
		metadataBucket := tx.Bucket([]byte(EVENT_METADATA_BUCKET))
		aggregateBucket := tx.Bucket([]byte(AGGREGATES_BUCKET))
		for _, event := range events {

			if err := aggregateBucket.Put([]byte(aggregateId+"::"+string(event.AggregateVersion)), []byte(event.Id)); err != nil {
				return err
			}

			if e, err := event.Serialize(); err != nil || eventsBucket.Put([]byte(event.Id), e) != nil {
				return err
			}

			offset, _ := metadataBucket.NextSequence()
			if em, err := event.EventMetaData.Serialize(); err != nil || metadataBucket.Put([]byte(strconv.FormatUint(offset, 10)), em) != nil {
				return err
			}
		}
		return nil
	})
}

func (store *BoltEventStore) GetEventMetaDataFrom(offset, count int) []EventMetaData {
	min := []byte(strconv.Itoa(offset))
	max := []byte(strconv.Itoa(offset + count))

	metas := make([]EventMetaData, 0)
	store.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(EVENT_METADATA_BUCKET)).Cursor()
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			m := new(EventMetaData)
			if err := m.Deserialize(v); err != nil {
				glog.Error("Error deserializing event", err)
				return err
			}
			metas = append(metas, *m)
		}
		return nil
	})
	return metas
}

func NewBoltStore(path string) *BoltEventStore {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		glog.Fatal("Error while opening bolt db", err)
	}
	db.Update(func(tx *bolt.Tx) error {
		createBucket(tx, EVENTS_BUCKET)
		createBucket(tx, EVENT_METADATA_BUCKET)
		createBucket(tx, AGGREGATES_BUCKET)
		return nil
	})
	return &BoltEventStore{db}
}
