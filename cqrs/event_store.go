package cqrs

//Represents a backing store to save aggregate's events
type EventStore interface {
	GetEvents(aggregateId string) []Event
	GetEventMetaDataFrom(offset, count int) []EventMetaData
	GetEvent(id string) Event
	SaveEvents(aggregateId string, events []Event) error
}
