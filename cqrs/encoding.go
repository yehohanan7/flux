package cqrs

import (
	"bytes"
	"encoding/gob"

	"github.com/golang/glog"
)

func deserialize(data []byte, target interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	return d.Decode(target)
}

func serialize(target interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(target)
	if err != nil {
		glog.Error("could not serialize data %v", target)
		return nil, err
	}
	return buffer.Bytes(), nil
}
