package mongodb

import (
	"time"

	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

type mongoAggregateRecord struct {
	Id      string                      `bson:"_id"`
	Version int                         `bson:"version"`
	Events  []mongoAggregateEventRecord `bson:"events"`
}

type mongoAggregateEventRecord struct {
	Id        string    `bson:"_id"`
	OccuredAt time.Time `bson:"occured_at"`
	Data      []byte    `bson:"data"`
}

type mongoEventRecord struct {
	Id          string    `bson:"_id"`
	AggregateId string    `bson:"aggregate_id"`
	OccuredAt   time.Time `bson:"occured_at"`
	Data        []byte    `bson:"data"`
}

type MongoEventStoreOptions struct {
	Session               *mgo.Session
	Database              string
	EventCollection       string
	AggregateCollection   string
	TransactionCollection string
}

func DefaultMongoEventStoreOptions() *MongoEventStoreOptions {
	return &MongoEventStoreOptions{
		Database:              "",
		EventCollection:       "event",
		AggregateCollection:   "aggregate",
		TransactionCollection: "transaction",
	}
}

// MongoEventStore implementation of the event store
type MongoEventStore struct {
	options *MongoEventStoreOptions
}

func (store *MongoEventStore) getEventCollection() *mgo.Collection {
	return store.options.Session.DB(store.options.Database).C(store.options.EventCollection)
}

func (store *MongoEventStore) getTransactionCollection() *mgo.Collection {
	return store.options.Session.DB(store.options.Database).C(store.options.TransactionCollection)
}

func (store *MongoEventStore) getAggregateCollection() *mgo.Collection {
	return store.options.Session.DB(store.options.Database).C(store.options.AggregateCollection)
}

func (store *MongoEventStore) GetEvents(aggregateId string) []Event {
	collection := store.getAggregateCollection()
	record := &mongoAggregateRecord{}
	err := collection.FindId(aggregateId).One(record)
	if err != nil {
		glog.Error("error while getting events ", err)
		return []Event{}
	}
	events := make([]Event, len(record.Events))
	for index, eventRecord := range record.Events {
		event := Event{}
		event.Deserialize(eventRecord.Data)
		events[index] = event
	}
	return events
}

func (store *MongoEventStore) GetEventMetaDataFrom(offset, count int) []EventMetaData {
	collection := store.getEventCollection()
	iter := collection.Find(nil).Skip(offset).Limit(count).Iter()
	record := &mongoEventRecord{}
	events := make([]EventMetaData, 0)
	for iter.Next(record) {
		meta := EventMetaData{}
		meta.Deserialize(record.Data)
		events = append(events, meta)
	}
	iter.Close()
	return events
}

func (store *MongoEventStore) createAggregate(aggregateId string) (bool, error) {
	aggregateCollection := store.getAggregateCollection()
	pipe := aggregateCollection.Pipe([]bson.M{
		bson.M{"$match": bson.M{"_id": aggregateId}},
		bson.M{"$project": bson.M{"version": 1}},
	})

	aggregateMetadata := make(map[string]interface{})
	err := pipe.One(aggregateMetadata)

	if err == mgo.ErrNotFound {
		err = aggregateCollection.Insert(mongoAggregateRecord{
			Id:      aggregateId,
			Version: 0,
		})
		return true, err
	}

	return false, err
}

func (store *MongoEventStore) SaveEvents(aggregateId string, events []Event) error {
	l := len(events)
	if l == 0 {
		return nil
	}

	assert := bson.M{"version": bson.M{"$eq": events[0].AggregateVersion - 1}}

	created, err := store.createAggregate(aggregateId)

	if err != nil {
		return err
	}

	if created {
		assert = nil
	}

	runner := txn.NewRunner(store.getTransactionCollection())
	ops := make([]txn.Op, l+2)

	serializedEvents := make([]mongoAggregateEventRecord, len(events))
	for index, event := range events {
		eventTime, err := time.Parse(time.ANSIC, event.OccuredAt)
		if err != nil {
			return err
		}
		serializedEvents[index] = mongoAggregateEventRecord{
			Id:        event.Id,
			OccuredAt: eventTime,
			Data:      event.Serialize(),
		}

		ops[index+2] = txn.Op{
			C:  store.options.EventCollection,
			Id: event.Id,
			Insert: mongoEventRecord{
				Id:          event.Id,
				AggregateId: aggregateId,
				OccuredAt:   eventTime,
				Data:        event.EventMetaData.Serialize(),
			},
		}
	}

	ops[0] = txn.Op{
		C:      store.options.AggregateCollection,
		Id:     aggregateId,
		Assert: assert,
		Update: bson.M{
			"$push": bson.M{
				"events": bson.M{
					"$each": serializedEvents,
				},
			},
		},
	}

	ops[1] = txn.Op{
		C:  store.options.AggregateCollection,
		Id: aggregateId,
		Update: bson.M{
			"$set": bson.M{
				"version": events[l-1].AggregateVersion,
			},
		},
	}

	id := bson.NewObjectId()
	return runner.Run(ops, id, nil)
}

func (store *MongoEventStore) GetEvent(id string) Event {
	eventCollection := store.getEventCollection()
	eventRecord := &mongoEventRecord{}
	err := eventCollection.FindId(id).One(eventRecord)
	if err != nil {
		glog.Fatal("Error while finding event ", err)
		return Event{}
	}
	aggregateCollection := store.getAggregateCollection()
	pipe := aggregateCollection.Pipe([]bson.M{
		bson.M{"$unwind": "$events"},
		bson.M{"$match": bson.M{"events._id": eventRecord.Id}},
		bson.M{"$project": bson.M{"events": []string{"$events"}}},
	})
	aggregateRecord := mongoAggregateRecord{}
	err = pipe.One(&aggregateRecord)
	if err != nil {
		glog.Error("Error while getting event ", err)
		return Event{}
	}
	aggregateEventRecord := aggregateRecord.Events[0]
	event := Event{}
	event.Deserialize(aggregateEventRecord.Data)
	return event
}

func NewEventStore(options *MongoEventStoreOptions) *MongoEventStore {
	return &MongoEventStore{options: options}
}
