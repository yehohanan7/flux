package cqrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var stores = []EventStore{NewInMemoryEventStore(), NewBoltEventStore("/tmp/temp.db")}

type EventPayload struct {
	Data string
}

func events() []Event {
	e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
	e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})
	return []Event{e1, e2}
}

func TestSaveEvents(t *testing.T) {
	for _, store := range stores {
		for _, aggregateId := range []string{"a-id", "b-id"} {
			err := store.SaveEvents(aggregateId, events())

			assert.Nil(t, err)
			assert.Len(t, store.GetEvents(aggregateId), 2)
			assert.Equal(t, "payload", store.GetEvents(aggregateId)[0].Payload.(EventPayload).Data)
		}
		assert.Len(t, store.GetAllEvents(), 4)
	}
}

func TestGetEvent(t *testing.T) {
	for _, store := range stores {
		expected := NewEvent("aid", 0, "payload")
		err := store.SaveEvents("a-id", []Event{expected})

		actual := store.GetEvent(expected.Id)

		assert.Nil(t, err)
		assert.Equal(t, expected.Id, actual.Id)
	}

}
