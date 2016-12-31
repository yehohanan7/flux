package cqrs

import (
	"encoding/gob"
	"reflect"

	"github.com/golang/glog"
)

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
	if err != nil {
		glog.Warningf("error while saving events for aggregate %v", aggregate)
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
	hm := buildHandlerMap(entity)
	for eventType, _ := range hm {
		gob.Register(reflect.New(eventType))
	}

	aggregate := Aggregate{
		Id:       id,
		version:  0,
		events:   []Event{},
		entity:   entity,
		handlers: hm,
		store:    store,
		name:     reflect.TypeOf(entity).String(),
	}

	aggregate.Apply(store.GetEvents(id)...)
	return aggregate
}
