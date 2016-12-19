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
	err := aggregate.store.SaveEvents(aggregate.Id, aggregate.events)
	if err == nil {
		aggregate.events = []Event{}
	}
	return err
}

func (aggregate *Aggregate) Update(payloads ...interface{}) {
	for _, payload := range payloads {
		event := NewEvent(aggregate.name, aggregate.version, payload)
		aggregate.events = append(aggregate.events, event)
		aggregate.Apply(event)
	}
}

func (aggregate *Aggregate) Apply(events ...Event) {
	for _, e := range events {
		payload := e.Payload
		if handler, ok := aggregate.handlers[reflect.TypeOf(payload)]; ok {
			handler(aggregate.entity, payload)
			aggregate.version++
		}
	}
}

func NewAggregate(id string, entity interface{}, store EventStore) Aggregate {
	aggregate := Aggregate{
		Id:       id,
		version:  0,
		events:   []Event{},
		entity:   entity,
		handlers: buildHandlerMap(entity),
		store:    store,
		name:     reflect.TypeOf(entity).Name(),
	}
	aggregate.Apply(store.GetEvents(id)...)
	return aggregate
}
