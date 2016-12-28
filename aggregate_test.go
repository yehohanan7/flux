package cqrs

import "testing"

import "github.com/stretchr/testify/assert"

type TestEvent struct {
}

type TestEntity struct {
	Aggregate
	handled bool
}

func (entity *TestEntity) HandleEvent(event TestEvent) {
	entity.handled = true
}

func newEntityWithAggregate() *TestEntity {
	entity := new(TestEntity)
	entity.Aggregate = NewAggregate("aggregate-id", entity, NewInMemoryEventStore())
	return entity
}

func TestDefaultVersion(t *testing.T) {
	entity := newEntityWithAggregate()

	assert.Equal(t, 0, entity.version)
}

func TestEventHandling(t *testing.T) {
	entity := newEntityWithAggregate()

	entity.Update(TestEvent{})

	assert.True(t, entity.handled)
}

func TestUpdateAggregateVersion(t *testing.T) {
	entity := newEntityWithAggregate()

	entity.Update(TestEvent{}, TestEvent{})

	assert.Equal(t, 2, entity.version)
}

func TestEventAggregateVersion(t *testing.T) {
	entity := newEntityWithAggregate()

	entity.Update(TestEvent{}, TestEvent{})

	assert.Equal(t, 0, entity.events[0].AggregateVersion)
	assert.Equal(t, 1, entity.events[1].AggregateVersion)
}

func TestUnknownEvent(t *testing.T) {
	entity := newEntityWithAggregate()

	assert.NotPanics(t, func() { entity.Update("unknown string event") })

	assert.False(t, entity.handled)
}

func TestAggregateSaveEvents(t *testing.T) {
	store := NewInMemoryEventStore()
	entity := new(TestEntity)
	entity.Aggregate = NewAggregate("aggregate-id", entity, store)
	entity.Update(TestEvent{}, TestEvent{})

	entity.Save()

	assert.Len(t, store.GetEvents("aggregate-id"), 2)
	assert.Empty(t, entity.events)
}
