package cqrs

type EventStore interface {
	GetEvents(aggregateId string) []Event
	SaveEvents(aggregateId string, events []Event) error
}
