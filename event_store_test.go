package cqrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var stores = []EventStore{NewInMemoryEventStore(), NewBoltEventStore("/tmp/temp.db")}

func events() []Event {
	e1 := NewEvent("sample_aggregate", 1, "payload")
	e2 := NewEvent("sample_aggregate", 2, "payload")
	return []Event{e1, e2}
}

func TestSaveEvents(t *testing.T) {
	for _, store := range stores {

		err1 := store.SaveEvents("a-id", events())
		err2 := store.SaveEvents("b-id", events())

		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Len(t, store.GetEvents("a-id"), 2)
		assert.Len(t, store.GetEvents("b-id"), 2)
	}

}
