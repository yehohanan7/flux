package cqrs

import "reflect"

const defaultOffset = 0

//Consumes events from the command component
type EventConsumer struct {
	handlerClass interface{}
	handlers     handlermap
}

//Send event to the consumer
func (consumer *EventConsumer) send(event Event) {
	payload := event.Payload
	if handler, ok := consumer.handlers[reflect.TypeOf(payload)]; ok {
		handler(consumer.handlerClass, payload)
	}
}

func (consumer *EventConsumer) Start() error {
	return nil
}

func (consumer *EventConsumer) Stop() error {
	return nil
}

//Create new consumer
func NewConsumer(url string, handlerClass interface{}, store OffsetStore) *EventConsumer {
	return &EventConsumer{handlerClass, buildHandlerMap(handlerClass)}
}
