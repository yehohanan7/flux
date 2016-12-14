package cqrs

import (
	"reflect"

	uuid "github.com/satori/go.uuid"
)

type Event struct {
	Id               string
	Type             string
	Payload          interface{}
	AggregateVersion int
	AggregateName    string
}

func NewEvent(payload interface{}, aggregate *Aggregate) Event {
	return Event{
		Id:               uuid.NewV4().String(),
		Payload:          payload,
		AggregateVersion: aggregate.version,
		AggregateName:    reflect.TypeOf(aggregate.entity).Name(),
	}
}
