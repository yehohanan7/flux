package cqrs

type EventConsumer interface {
	Start(eventCh chan interface{}) error
}
