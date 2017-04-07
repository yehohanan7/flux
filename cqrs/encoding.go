package cqrs

import (
	"bytes"
	"encoding/gob"

	"github.com/golang/glog"
)

func deserialize(data []byte, target interface{}) {
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	if err := d.Decode(target); err != nil {
		glog.Fatal("error while decoding ", err)
	}
}

func serialize(target interface{}) []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(target)
	if err != nil {
		glog.Fatal("could not serialize data ", err)
	}
	return buffer.Bytes()
}
