package cqrs

import uuid "github.com/satori/go.uuid"

type Event struct {
	Id               string
	Type             string
	Payload          interface{}
	OccuredAt        string
	AggregateVersion int
	AggregateName    string
}

func NewEvent(aggregateName string, aggregateVersion int, payload interface{}) Event {
	return Event{
		Id:               uuid.NewV4().String(),
		Payload:          payload,
		AggregateVersion: aggregateVersion,
		AggregateName:    aggregateName,
	}
}
