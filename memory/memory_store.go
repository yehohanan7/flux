package memory

import (
	. "github.com/yehohanan7/flux/cqrs"
)

//InMemory implementation of the event store
type InMemoryEventStore struct {
	events     []Event            //collection of events
	aggregates map[string][]Event //map of aggregate and all events occured on it
}

func (store *InMemoryEventStore) GetEvents(aggregateId string) []Event {
	return store.aggregates[aggregateId]
}

func (store *InMemoryEventStore) GetAllEventsFrom(offset, count int) []Event {
	for i, _ := range store.events {
		if i == offset {
			if i+count > len(store.events) {
				return store.events[i:]
			} else {
				return store.events[i : i+count]
			}

		}
	}
	return []Event{}
}

func (store *InMemoryEventStore) SaveEvents(aggregateId string, events []Event) error {
	if _, ok := store.aggregates[aggregateId]; !ok {
		store.aggregates[aggregateId] = make([]Event, 0)
	}
	store.aggregates[aggregateId] = append(store.aggregates[aggregateId], events...)
	store.events = append(store.events, events...)
	return nil
}

func (store *InMemoryEventStore) GetEvent(id string) Event {
	var event Event
	for _, e := range store.events {
		if e.Id == id {
			return e
		}
	}
	return event
}

func NewEventStore() EventStore {
	return &InMemoryEventStore{make([]Event, 0), make(map[string][]Event)}
}
