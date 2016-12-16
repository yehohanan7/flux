package cqrs

import "reflect"

type Aggregate struct {
	Id       string
	name     string
	version  int
	events   []Event
	entity   interface{}
	handlers HandlerMap
	store    EventStore
}

func (aggregate *Aggregate) Save() error {
	return aggregate.store.SaveEvents(aggregate.Id, aggregate.events)
}

func (aggregate *Aggregate) Update(events ...interface{}) {
	for _, event := range events {
		aggregate.events = append(aggregate.events, NewEvent(aggregate.name, aggregate.version, event))
		aggregate.Apply(event)
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

func NewAggregate(id string, entity interface{}, store EventStore) Aggregate {
	return Aggregate{
		Id:       id,
		version:  0,
		events:   []Event{},
		entity:   entity,
		handlers: buildHandlerMap(entity),
		store:    store,
		name:     reflect.TypeOf(entity).Name(),
	}
}
