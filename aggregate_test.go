package cqrs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
	entity.Aggregate = NewAggregate("aggregate-id", entity, NewEventStore())
	return entity
}

var _ = Describe("Aggregare", func() {

	var (
		aggregateId string     = "aggregate-id"
		store       EventStore = NewEventStore()
		entity      *TestEntity
	)

	BeforeEach(func() {
		entity = new(TestEntity)
		entity.Aggregate = NewAggregate(aggregateId, entity, store)
	})

	_ = Describe("Creating New Aggregate", func() {
		It("should set version to 0", func() {
			Expect(entity.version).To(Equal(0))
		})
	})

	_ = Describe("Executing update on aggregate", func() {
		It("should handle the event", func() {
			entity.Update(TestEvent{})

			Expect(entity.handled).To(BeTrue())
		})

		It("should update aggregate version", func() {
			entity.Update(TestEvent{}, TestEvent{})

			Expect(entity.version).To(Equal(2))
		})

		It("should update event's aggregate version", func() {
			entity.Update(TestEvent{}, TestEvent{})

			Expect(entity.events[0].AggregateVersion).To(Equal(0))
			Expect(entity.events[1].AggregateVersion).To(Equal(1))
		})

		It("should not panic for unknown events", func() {
			Expect(func() { entity.Update("unknown string event") }).ShouldNot(Panic())

			Expect(entity.handled).To(BeFalse())
		})
	})

	_ = Describe("Saving an aggregate", func() {
		It("Should store the events and clear state", func() {
			entity.Update(TestEvent{}, TestEvent{})

			entity.Save()

			Expect(len(entity.events)).To(Equal(0))
			Expect(len(store.GetEvents(aggregateId))).To(Equal(2))
		})
	})

})
