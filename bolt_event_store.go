package cqrs

type BoltEventStore struct {
}

func (store *BoltEventStore) GetEvents(aggregateId string) []Event {
	return nil
}

func (store *BoltEventStore) SaveEvents(aggregateId string, events []Event) error {
	return nil
}

func NewBoltEventStore() EventStore {
	return &BoltEventStore{}
}
