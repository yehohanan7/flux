package cqrs

import (
	"reflect"

	"github.com/golang/glog"
)

//Aggregate as in the DDD world
type Aggregate struct {
	Id       string
	name     string
	Version  int
	Events   []Event
	entity   interface{}
	handlers Handlers
	store    EventStore
}

//Save the events accumulated so far
func (aggregate *Aggregate) Save() error {
	err := aggregate.store.SaveEvents(aggregate.Id, aggregate.Events)
	if err != nil {
		glog.Warning("error while saving events for aggregate ", aggregate, err)
		return err
	}
	aggregate.Events = []Event{}
	return err
}

//Update the event
func (aggregate *Aggregate) Update(payloads ...interface{}) {
	for _, payload := range payloads {
		event := NewEvent(aggregate.name, aggregate.Version, payload)
		aggregate.Events = append(aggregate.Events, event)
		aggregate.apply(event)
	}
}

//Apply events
func (aggregate *Aggregate) apply(events ...Event) {
	for _, e := range events {
		payload := e.Payload
		if handler, ok := aggregate.handlers[reflect.TypeOf(payload)]; ok {
			handler(aggregate.entity, payload)
			aggregate.Version++
		}
	}
}

//Create new aggregate with a backing event store
func NewAggregate(id string, entity interface{}, store EventStore) Aggregate {
	aggregate := Aggregate{
		Id:       id,
		Version:  0,
		Events:   []Event{},
		entity:   entity,
		handlers: NewHandlers(entity),
		store:    store,
		name:     reflect.TypeOf(entity).String(),
	}
	return aggregate
}

//Get the aggregate
func GetAggregate(id string, entity interface{}, store EventStore) Aggregate {
	aggregate := NewAggregate(id, entity, store)
	aggregate.apply(store.GetEvents(id)...)
	return aggregate
}
