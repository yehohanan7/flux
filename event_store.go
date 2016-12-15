package cqrs

type EventStore interface {
	GetEvents(aggregateId, aggregateType string) []Event
	SaveEvents(aggregateId, events []Event) error
}
