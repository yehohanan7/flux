package boltdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yehohanan7/cqrs"
)

const DB_PATH string = "/tmp/temp.db"

func init() {
	os.Remove(DB_PATH)
}

type EventPayload struct {
	Data string
}

func events() []cqrs.Event {
	e1 := cqrs.NewEvent("sample_aggregate", 1, EventPayload{"payload"})
	e2 := cqrs.NewEvent("sample_aggregate", 2, EventPayload{"payload"})
	return []cqrs.Event{e1, e2}
}

func TestEventKey(t *testing.T) {
	key := buildKey("aggregateId", cqrs.NewEvent("AggregateName", 1, "payload"))

	assert.Equal(t, "aggregateId_1", key)
}

func TestSaveEvents(t *testing.T) {
	store := NewEventStore(DB_PATH)
	for _, aggregateId := range []string{"a-id", "b-id"} {
		err := store.SaveEvents(aggregateId, events())

		assert.Nil(t, err)
		assert.Len(t, store.GetEvents(aggregateId), 2)
		assert.Equal(t, "payload", store.GetEvents(aggregateId)[0].Payload.(EventPayload).Data)
	}
	assert.Len(t, store.GetAllEvents(), 4)
}

func TestGetEvent(t *testing.T) {
	store := NewEventStore(DB_PATH)
	expected := cqrs.NewEvent("aid", 0, EventPayload{"payload"})
	err := store.SaveEvents("a-id", []cqrs.Event{expected})

	actual := store.GetEvent(expected.Id)

	assert.Nil(t, err)
	assert.Equal(t, expected.Id, actual.Id)
}
