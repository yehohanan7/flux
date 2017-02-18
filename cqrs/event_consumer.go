package cqrs

type EventConsumer interface {
	Start() error
	Stop() error
}
