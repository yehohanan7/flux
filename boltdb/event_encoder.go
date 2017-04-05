package boltdb

import (
	"bytes"
	"encoding/gob"

	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
)

func encodeEvent(event Event) []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(event)
	if err != nil {
		glog.Fatal("could not serialize event %v", event)
	}
	return buffer.Bytes()
}

func decodeEvent(data []byte) Event {
	e := new(Event)
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err := d.Decode(e)
	if err != nil {
		glog.Fatal("could not deserialize event %v", e)
	}
	return *e
}
