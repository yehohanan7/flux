package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

func createBucket(tx *bolt.Tx, name string) *bolt.Bucket {
	b, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		glog.Fatal("Error while initializing db", err)
	}
	return b
}
