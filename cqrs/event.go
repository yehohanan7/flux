package cqrs

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

//Every action on an aggregate emits an event, which is wrapped and saved
type Event struct {
	Id               string      `json:"id"`
	OccuredAt        string      `json:"occured_at"`
	AggregateVersion int         `json:"aggregate_version"`
	AggregateName    string      `json:"aggregate_name"`
	Payload          interface{} `json:"payload"`
}

//Create new event
func NewEvent(aggregateName string, aggregateVersion int, payload interface{}) Event {
	return Event{
		Id:               uuid.NewV4().String(),
		AggregateVersion: aggregateVersion,
		AggregateName:    aggregateName,
		OccuredAt:        time.Now().Format(time.ANSIC),
		Payload:          payload,
	}
}
