package boltdb

import (
	"bytes"
	"encoding/gob"

	"github.com/golang/glog"
)

func encodeEvent(event Event) []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(e)
	if err != nil {
		glog.Errorf("could not serialize event %v", e)
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
		glog.Errorf("could not deserialize event %v", e)
	}
	return e
}
