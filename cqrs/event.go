package cqrs

import (
	"encoding/gob"

	"time"

	uuid "github.com/satori/go.uuid"
)

//Every action on an aggregate emits an event, which is wrapped and saved
type Event struct {
	Id               string      `json:"id"`
	Payload          interface{} `json:"payload"`
	OccuredAt        string      `json:"occured_at"`
	AggregateVersion int         `json:"aggregate_version"`
	AggregateName    string      `json:"aggregate_name"`
}

//Create new event
func NewEvent(aggregateName string, aggregateVersion int, payload interface{}) Event {
	gob.Register(payload)
	return Event{
		Id:               uuid.NewV4().String(),
		Payload:          payload,
		AggregateVersion: aggregateVersion,
		AggregateName:    aggregateName,
		OccuredAt:        time.Now().Format(time.ANSIC),
	}
}
