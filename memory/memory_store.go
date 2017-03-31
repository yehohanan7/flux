package memory

import . "github.com/yehohanan7/flux/cqrs"

//InMemory implementation of the event store
type InMemoryEventStore struct {
	events       []Event
	eventMap     map[string]Event
	aggregates   map[string][]Event
	aggregateMap map[string]Aggregate
}

func (store *InMemoryEventStore) GetEvents(aggregateId string) []Event {
	return store.aggregates[aggregateId]
}

func asMetaData(events []Event) []EventMetaData {
	m := make([]EventMetaData, len(events))
	for i := 0; i < len(events); i++ {
		m[i] = events[i].EventMetaData
	}
	return m
}

func (store *InMemoryEventStore) GetEventMetaDataFrom(offset, count int) []EventMetaData {
	for i, _ := range store.events {
		if i == offset {
			if i+count > len(store.events) {
				return asMetaData(store.events[i:])
			} else {
				return asMetaData(store.events[i : i+count])
			}
		}
	}
	return []EventMetaData{}
}

func (store *InMemoryEventStore) SaveEvents(aggregateId string, events []Event) error {
	if _, ok := store.aggregates[aggregateId]; !ok {
		store.aggregates[aggregateId] = make([]Event, 0)
	}
	store.aggregates[aggregateId] = append(store.aggregates[aggregateId], events...)

	for _, e := range events {
		store.events = append(store.events, e)
		store.eventMap[e.Id] = e
	}

	return nil
}

func (store *InMemoryEventStore) GetEvent(id string) Event {
	return store.eventMap[id]
}

func (store *InMemoryEventStore) SaveAggregate(aggregate Aggregate) error {
	store.aggregateMap[aggregate.Id] = aggregate
	return nil
}

func (store *InMemoryEventStore) GetAggregate(id string) Aggregate {
	return store.aggregateMap[id]
}

func NewEventStore() EventStore {
	return &InMemoryEventStore{make([]Event, 0), make(map[string]Event), make(map[string][]Event), make(map[string]Aggregate)}
}
