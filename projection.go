package cqrs

import (
	"reflect"

	. "github.com/yehohanan7/cqrs/event"
)

type Projection struct {
	entity   interface{}
	handlers HandlerMap
}

func (projection *Projection) Apply(events []Event) {
	for _, e := range events {
		payload := e.Payload
		if handler, ok := projection.handlers[reflect.TypeOf(payload)]; ok {
			handler(projection.entity, payload)
		}
	}
}

func NewProjection(aggregateId string, model interface{}, store EventStore) Projection {
	projection := Projection{model, buildHandlerMap(model)}
	projection.Apply(store.GetEvents(aggregateId))
	return projection
}
