package cqrs

type EventStore interface {
	GetEvents(aggregateId string) []Event
	SaveEvents(aggregateId string, events []Event) error
}

type InMemoryEventStore struct {
	events map[string][]Event
}

func (store *InMemoryEventStore) GetEvents(aggregateId string) []Event {
	return nil
}

func (store *InMemoryEventStore) SaveEvents(aggregateId string, events []Event) error {
	return nil
}

func NewInMemoryEventStore() EventStore {
	return &InMemoryEventStore{make(map[string][]Event)}
}
