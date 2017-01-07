package cqrs

import (
	"encoding/gob"
	"reflect"

	"github.com/golang/glog"
)

//Aggregate as in the DDD world
type Aggregate struct {
	Id       string
	name     string
	version  int
	events   []Event
	entity   interface{}
	handlers handlermap
	store    EventStore
}

//Save the events acculated so far
func (aggregate *Aggregate) Save() error {
	err := aggregate.store.SaveEvents(aggregate.Id, aggregate.events)
	if err != nil {
		glog.Warning("error while saving events for aggregate ", aggregate, err)
	}
	aggregate.events = []Event{}
	return err
}

//Update the event
func (aggregate *Aggregate) Update(payloads ...interface{}) {
	for _, payload := range payloads {
		event := NewEvent(aggregate.name, aggregate.version, payload)
		aggregate.events = append(aggregate.events, event)
		aggregate.Apply(event)
	}
}

//Apply events
func (aggregate *Aggregate) Apply(events ...Event) {
	for _, e := range events {
		payload := e.Payload
		if handler, ok := aggregate.handlers[reflect.TypeOf(payload)]; ok {
			handler(aggregate.entity, payload)
			aggregate.version++
		}
	}
}

//Create new aggregate with a backing event store
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
