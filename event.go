package cqrs

import (
	"bytes"

	"encoding/gob"

	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

type Event struct {
	Id               string
	Payload          interface{}
	OccuredAt        string
	AggregateVersion int
	AggregateName    string
}

func (e *Event) Deserialize(data []byte) {
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err := d.Decode(e)
	if err != nil {
		glog.Errorf("could not deserialize event %v", e)
	}
}

func (e *Event) Serialize() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(e)
	if err != nil {
		glog.Errorf("could not serialize event %v", e)
	}
	return buffer.Bytes()
}

func NewEvent(aggregateName string, aggregateVersion int, payload interface{}) Event {
	gob.Register(payload)
	return Event{
		Id:               uuid.NewV4().String(),
		Payload:          payload,
		AggregateVersion: aggregateVersion,
		AggregateName:    aggregateName,
	}
}
