package cqrs

import uuid "github.com/satori/go.uuid"

type Event struct {
	Id               string      `json:"id"`
	Type             string      `json:"type"`
	Payload          interface{} `json:"payload"`
	OccuredAt        string      `json:"occured_at"`
	AggregateVersion int         `json:"aggregate_version"`
	AggregateName    string      `json:"aggregate_name"`
}

func NewEvent(aggregateName string, aggregateVersion int, payload interface{}) Event {
	return Event{
		Id:               uuid.NewV4().String(),
		Payload:          payload,
		AggregateVersion: aggregateVersion,
		AggregateName:    aggregateName,
	}
}
