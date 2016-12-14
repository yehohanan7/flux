package cqrs

import "testing"

import "github.com/stretchr/testify/assert"

func TestNewEvent(t *testing.T) {
	event := NewEvent("some payload", &Aggregate{entity: TestEntity{}, version: 1})

	assert.NotEqual(t, "", event.Id)
	assert.True(t, len(event.Id) > 25)
	assert.Equal(t, "TestEntity", event.AggregateName)
	assert.Equal(t, 1, event.AggregateVersion)
}
