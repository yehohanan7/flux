package boltdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
			event.Deserialize(v)
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
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			events = append(events, store.GetEvent(string(v)))
		}
		return nil
	})
	return events
}

func aggregateKey(id string, version int) []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte(id))
	v := uint16(version)
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		glog.Fatal(fmt.Sprintf("error while forming key for aggregae %s & version %d\n", id, version))
	}
	return buf.Bytes()
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []Event) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		eb := tx.Bucket([]byte(EVENTS_BUCKET))
		mb := tx.Bucket([]byte(EVENT_METADATA_BUCKET))
		ab := tx.Bucket([]byte(AGGREGATES_BUCKET))

		lastKey := aggregateKey(aggregateId, events[0].AggregateVersion-1)
		newKey := aggregateKey(aggregateId, events[0].AggregateVersion)

		if ab.Get(newKey) != nil || (events[0].AggregateVersion != 0 && ab.Get(lastKey) == nil) {
			return Conflict
		}

		for _, event := range events {

			if err := ab.Put(aggregateKey(aggregateId, event.AggregateVersion), []byte(event.Id)); err != nil {
				return err
			}

			bytes := event.Serialize()
			if err := eb.Put([]byte(event.Id), bytes); err != nil {
				return err
			}

			offset, _ := mb.NextSequence()
			bytes = event.EventMetaData.Serialize()
			if err := mb.Put([]byte(strconv.FormatUint(offset-uint64(1), 10)), bytes); err != nil {
				return err
			}
		}
		return nil
	})
}

func (store *BoltEventStore) GetEventMetaDataFrom(offset, count int) []EventMetaData {

	metas := make([]EventMetaData, 0)
	store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EVENT_METADATA_BUCKET))
		for i := offset; i < offset+count; i++ {
			v := b.Get([]byte(strconv.Itoa(i)))
			if v != nil && len(v) > 0 {
				m := new(EventMetaData)
				m.Deserialize(v)
				metas = append(metas, *m)
			}
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
