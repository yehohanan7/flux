package cqrs

import "reflect"

type Aggregate struct {
	Id       string
	version  int
	events   []Event
	entity   interface{}
	handlers HandlerMap
}

func (aggregate *Aggregate) Update(events ...interface{}) {
	for _, e := range events {
		aggregate.events = append(aggregate.events, NewEvent(e, aggregate))
		aggregate.Apply(e)
	}
}

func (aggregate *Aggregate) Apply(events ...interface{}) {
	for _, e := range events {
		if handler, ok := aggregate.handlers[reflect.TypeOf(e)]; ok {
			handler(aggregate.entity, e)
			aggregate.version++
		}
	}
}

func NewAggregate(id string, entity interface{}) Aggregate {
	return Aggregate{id, 0, []Event{}, entity, buildHandlerMap(entity)}
}
