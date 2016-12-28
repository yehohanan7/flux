package cqrs

type EventStore interface {
	GetAllEvents() []Event
	GetEvents(aggregateId string) []Event
	SaveEvents(aggregateId string, events []Event) error
}
