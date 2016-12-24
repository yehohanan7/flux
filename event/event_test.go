package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SomeData struct {
	Data string
}

type SomePayLoad struct {
	SomeData SomeData
}

func TestNewEvent(t *testing.T) {
	event := NewEvent("SomeAggregate", 1, SomePayLoad{SomeData{"some data"}})

	assert.True(t, len(event.Id) > 25)
	assert.Equal(t, "SomeAggregate", event.AggregateName)
	assert.Equal(t, 1, event.AggregateVersion)
}

func TestSerialize(t *testing.T) {
	e := NewEvent("SomeAggregate", 1, SomePayLoad{SomeData{"some data"}})

	bytes := e.Serialize()

	event := new(Event)
	event.Deserialize(bytes)
	assert.Equal(t, "SomeAggregate", event.AggregateName)
	assert.Equal(t, "some data", event.Payload.(SomePayLoad).SomeData.Data)
}
