package cqrs

import "testing"

import "github.com/stretchr/testify/assert"

func TestNewEvent(t *testing.T) {
	event := NewEvent("SomeAggregate", 1, "some payload")

	assert.True(t, len(event.Id) > 25)
	assert.Equal(t, "SomeAggregate", event.AggregateName)
	assert.Equal(t, 1, event.AggregateVersion)
}
