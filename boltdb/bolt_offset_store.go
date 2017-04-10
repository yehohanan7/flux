package boltdb

import (
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/yehohanan7/flux/cqrs"
)

const OFFSET_BUCKET = "OFFSETS"

type BoltOffsetStore struct {
	db *bolt.DB
}

func (store *BoltOffsetStore) SaveOffset(value int) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(OFFSET_BUCKET))
		return b.Put([]byte("offset"), []byte(strconv.Itoa(value)))
	})
}

func (store *BoltOffsetStore) GetLastOffset() (int, error) {
	offset := -1
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(OFFSET_BUCKET))
		v := b.Get([]byte("offset"))
		o, err := strconv.Atoi(string(v))
		offset = o
		return err
	})
	return offset, err
}

//New bolt offset store
func NewOffsetStore(path string) cqrs.OffsetStore {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		glog.Fatal("Error while opening bolt db", err)
	}
	db.Update(func(tx *bolt.Tx) error {
		createBucket(tx, OFFSET_BUCKET)
		return nil
	})
	return &BoltOffsetStore{db}
}
