package cqrs

import "reflect"

type Aggregate struct {
	Id       string
	name     string
	version  int
	events   []Event
	entity   interface{}
	handlers HandlerMap
	repo     AggregateRepository
}

func (aggregate *Aggregate) Save() error {
	return aggregate.repo.Save(*aggregate)
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

func NewAggregate(id string, entity interface{}, repo AggregateRepository) Aggregate {
	return Aggregate{
		Id:       id,
		version:  0,
		events:   []Event{},
		entity:   entity,
		handlers: buildHandlerMap(entity),
		repo:     repo,
		name:     reflect.TypeOf(entity).Name(),
	}
}
