package boltdb

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

func createBucket(tx *bolt.Tx, name string) {
	_, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		glog.Fatal("Error while initializing db", err)
	}
}

func save(bucket *bolt.Bucket, key []byte, data interface{}) error {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		glog.Error("could not serialize data %v", data)
		return err
	}
	return bucket.Put(key, buffer.Bytes())
}
