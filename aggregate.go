package cqrs

import "reflect"

type Aggregate struct {
	Id       string
	version  int
	events   []Event
	entity   interface{}
	handlers HandlerMap
}

func (aggregate *Aggregate) Update(payload interface{}) {
	event := NewEvent(payload, aggregate)
	if handler, ok := aggregate.handlers[reflect.TypeOf(payload)]; ok {
		aggregate.events = append(aggregate.events, event)
		handler(aggregate.entity, payload)
		aggregate.version++
	}
}

func NewAggregate(id string, entity interface{}) Aggregate {
	return Aggregate{id, 0, []Event{}, entity, buildHandlerMap(entity)}
}
