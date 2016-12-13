package cqrs

import "reflect"

type Aggregate struct {
	events   []Event
	entity   interface{}
	handlers HandlerMap
}

func (aggregate *Aggregate) Update(payload interface{}) {
	event := Event{Payload: payload}
	if handler, ok := aggregate.handlers[reflect.TypeOf(payload)]; ok {
		aggregate.events = append(aggregate.events, event)
		handler(aggregate.entity, payload)
	}
}

func NewAggregate(entity interface{}) Aggregate {
	return Aggregate{[]Event{}, entity, buildHandlerMap(entity)}
}
