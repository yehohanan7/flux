package cqrs

import "testing"

import "github.com/stretchr/testify/assert"

type TestEvent struct {
}

type TestEntity struct {
	Aggregate
	handled bool
}

func (entity *TestEntity) Handle(event TestEvent) {
	entity.handled = true
}

func TestEventHandling(t *testing.T) {
	entity := new(TestEntity)
	entity.Aggregate = NewAggregate("aggregate-id", entity)

	entity.Update(TestEvent{})
	entity.Update(TestEvent{})

	assert.True(t, entity.handled)
	assert.Equal(t, 0, entity.events[0].AggregateVersion)
	assert.Equal(t, 1, entity.events[1].AggregateVersion)
	assert.Equal(t, 2, entity.version)
}

func TestUnknownEvent(t *testing.T) {
	entity := new(TestEntity)
	entity.Aggregate = NewAggregate("some-id", entity)

	assert.NotPanics(t, func() { entity.Update("unknown string event") })
	assert.False(t, entity.handled)
}

func TestDefaultVersion(t *testing.T) {
	entity := new(TestEntity)
	entity.Aggregate = NewAggregate("some-id", entity)

	assert.Equal(t, 0, entity.version)
}
