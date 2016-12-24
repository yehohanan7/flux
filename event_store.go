package cqrs

import . "github.com/yehohanan7/cqrs/event"

type EventStore interface {
	GetEvents(aggregateId string) []Event
	SaveEvents(aggregateId string, events []Event) error
}
