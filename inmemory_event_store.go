package cqrs

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

func (store *InMemoryEventStore) GetAllEvents() []Event {
	all := make([]Event, 0)
	for _, events := range store.events {
		all = append(all, events...)
	}
	return all
}

func NewInMemoryEventStore() EventStore {
	return &InMemoryEventStore{make(map[string][]Event)}
}