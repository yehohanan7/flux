package mongodb

import (
	"github.com/yehohanan7/flux/cqrs"
	mgo "gopkg.in/mgo.v2"
)

type mongoOffsetRecord struct {
	Id     string `bson:"_id"`
	Offset int    `bson:"offset"`
}

type MongoOffsetStoreOptions struct {
	Session    *mgo.Session
	Database   string
	Collection string
	StoreId    string
}

func DefaultMongoOffsetStoreOptions() *MongoOffsetStoreOptions {
	return &MongoOffsetStoreOptions{
		Database:   "",
		Collection: "offset",
	}
}

// MongoOffsetStore implementation of the offset store
type MongoOffsetStore struct {
	options *MongoOffsetStoreOptions
}

func (store *MongoOffsetStore) getCollection() *mgo.Collection {
	return store.options.Session.DB(store.options.Database).C(store.options.Collection)
}

func (store *MongoOffsetStore) SaveOffset(value int) error {
	collection := store.getCollection()
	_, err := collection.UpsertId(store.options.StoreId, mongoOffsetRecord{
		Id:     store.options.StoreId,
		Offset: value,
	})
	return err
}

func (store *MongoOffsetStore) GetLastOffset() (int, error) {
	collection := store.getCollection()
	query := collection.FindId(store.options.StoreId)
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
