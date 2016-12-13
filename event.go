package cqrs

type Event struct {
	Id      string
	Type    string
	Payload interface{}
}
