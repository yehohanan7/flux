package mongodb

import (
	"github.com/yehohanan7/flux/cqrs"
	mgo "gopkg.in/mgo.v2"
)

type mongoOffsetRecord struct {
	ID     string
	Offset int
}

type MongoOffsetStoreOptions struct {
	Session        *mgo.Session
	DatabaseName   string
	CollectionName string
	StoreID        string
}

func DefaultMongoOffsetStoreOptions() *MongoOffsetStoreOptions {
	return &MongoOffsetStoreOptions{
		DatabaseName:   "",
		CollectionName: "offset",
	}
}

type MongoOffsetStore struct {
	options *MongoOffsetStoreOptions
}

func (store *MongoOffsetStore) getCollection() *mgo.Collection {
	return store.options.Session.DB(store.options.DatabaseName).C(store.options.CollectionName)
}

func (store *MongoOffsetStore) SaveOffset(value int) error {
	collection := store.getCollection()
	_, err := collection.UpsertId(store.options.StoreID, mongoOffsetRecord{
		ID:     store.options.StoreID,
		Offset: value,
	})
	return err
}

func (store *MongoOffsetStore) GetLastOffset() (int, error) {
	collection := store.getCollection()
	query := collection.FindId(store.options.StoreID)
	offset := 0
	count, err := query.Count()
	if err == nil && count == 1 {
		record := mongoOffsetRecord{}
		err := query.One(&record)
		if err != nil {
			return 0, err
		}
		offset = record.Offset
	}
	return offset, err
}

func NewOffsetStore(options *MongoOffsetStoreOptions) cqrs.OffsetStore {
	return &MongoOffsetStore{options}
}
