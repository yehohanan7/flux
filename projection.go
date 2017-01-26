package cqrs

import (
	"reflect"
)

//Represents the read model in CQRS context
type Projection struct {
	entity   interface{}
	handlers handlermap
}

//Apply all events
func (projection *Projection) Apply(events []Event) {
	for _, e := range events {
		payload := e.Payload
		if handler, ok := projection.handlers[reflect.TypeOf(payload)]; ok {
			handler(projection.entity, payload)
		}
	}
}

//Create projection after applying events from store
func NewProjection(aggregateId string, model interface{}, store EventStore) Projection {
	projection := Projection{model, buildHandlerMap(model)}
	projection.Apply(store.GetEvents(aggregateId))
	return projection
}
