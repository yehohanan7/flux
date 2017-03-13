package cqrs

type EventConsumer interface {
	Start(eventCh, stopCh chan interface{}) error
}
