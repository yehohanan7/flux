package cqrs

import "reflect"

type Projection struct {
	entity   interface{}
	handlers HandlerMap
}

func (projection *Projection) Apply(events []interface{}) {
	for _, event := range events {
		if handler, ok := projection.handlers[reflect.TypeOf(event)]; ok {
			handler(projection.entity, event)
		}
	}
}

func NewProjection(model interface{}) Projection {
	return Projection{model, buildHandlerMap(model)}
}
