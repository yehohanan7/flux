package boltdb

import (
	. "github.com/yehohanan7/flux/cqrs"
)

//InMemory implementation of the event store
type BoltEventStore struct {
}

func (store *BoltEventStore) GetEvent(id string) Event {
	return Event{}
}

func (store *BoltEventStore) GetEvents(aggregateId string) []Event {
	return nil
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []Event) error {
	return nil
}

func (store *BoltEventStore) GetEventMetaDataFrom(offset, count int) []EventMetaData {
	return nil
}
