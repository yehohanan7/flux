package cqrs

//Represents a backing store to save aggregate's events
type EventStore interface {
	GetEvent(id string) Event
	GetEvents(aggregateId string) []Event
	SaveEvents(aggregateId string, events []Event) error
	GetEventMetaDataFrom(offset, count int) []EventMetaData
}
