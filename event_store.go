package cqrs

type EventStore interface {
	GetEvents(aggregateId string) []Event
	SaveEvents(aggregateId string, events []Event) error
}

type InMemoryEventStore struct {
	events map[string][]Event
}

func (store *InMemoryEventStore) GetEvents(aggregateId string) []Event {
	return store.events[aggregateId]
}

func (store *InMemoryEventStore) SaveEvents(aggregateId string, events []Event) error {
	if _, ok := store.events[aggregateId]; !ok {
		store.events[aggregateId] = make([]Event, 0)
	}
	store.events[aggregateId] = append(store.events[aggregateId], events...)
	return nil
}

func NewInMemoryEventStore() EventStore {
	return &InMemoryEventStore{make(map[string][]Event)}
}
