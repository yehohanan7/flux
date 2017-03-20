package cqrs

//Represents a backing store to save aggregate's events
type EventStore interface {
	GetEvents(aggregateId string) []Event
	GetAllEventsFrom(offset, count int) []Event
	GetEvent(id string) Event
	SaveEvents(aggregateId string, events []Event) error
}
