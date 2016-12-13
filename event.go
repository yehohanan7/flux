package cqrs

type Event struct {
	Id               string
	Type             string
	Payload          interface{}
	AggregateVersion int
}

func NewEvent(payload interface{}, aggregate *Aggregate) Event {
	return Event{Payload: payload, AggregateVersion: aggregate.version}
}
