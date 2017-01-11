package cqrs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var stores = []EventStore{NewEventStore()}

type EventPayload struct {
	Data string
}

func events() []Event {
	e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
	e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})
	return []Event{e1, e2}
}

var _ = Describe("InMemoryStore", func() {
	var store EventStore

	BeforeEach(func() {
		store = NewEventStore()
	})

	It("Should store events", func() {
		err1 := store.SaveEvents("aggregate-1", events())
		err2 := store.SaveEvents("aggregate-2", events())

		Expect(err1, err2).To(BeNil())
		Expect(store.GetEvents("aggregate-1")).To(HaveLen(2))
		Expect(store.GetEvents("aggregate-2")).To(HaveLen(2))
		Expect(store.GetAllEvents()).To(HaveLen(4))
	})

	var _ = Describe("Fetching all events from a secific event", func() {
		It("Should get the events", func() {
			e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
			e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})
			e3 := NewEvent("sample_aggregate", 3, EventPayload{"payload"})
			e4 := NewEvent("sample_aggregate", 4, EventPayload{"payload"})

			store.SaveEvents("aggregate1", []Event{e1, e2, e3, e4})

			events := store.GetAllEventsFrom(e2.Id, 2)
			Expect(events).To(HaveLen(2))
			Expect(events[0].Id).To(Equal(e2.Id))
			Expect(events[1].Id).To(Equal(e3.Id))
		})

		It("Should handle count gracefully", func() {
			e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
			e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})

			store.SaveEvents("aggregate1", []Event{e1, e2})

			events := store.GetAllEventsFrom(e2.Id, 5)

			Expect(events).To(HaveLen(1))
			Expect(events[0].Id).To(Equal(e2.Id))
		})
	})

	var _ = Describe("Fetching events of an aggregate", func() {
		It("Should get all events of an aggregate from a specific event", func() {
			e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
			e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})
			e3 := NewEvent("sample_aggregate", 3, EventPayload{"payload"})
			e4 := NewEvent("sample_aggregate", 4, EventPayload{"payload"})

			store.SaveEvents("aggregate1", []Event{e1, e2, e3})
			store.SaveEvents("aggregate2", []Event{e4})

			events := store.GetEventsFrom("aggregate1", e2.Id, 2)

			Expect(events).To(HaveLen(2))
			Expect(events[0].Id).To(Equal(e2.Id))
			Expect(events[1].Id).To(Equal(e3.Id))
		})

		It("Should handle large count value gracefully", func() {
			e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
			e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})

			store.SaveEvents("aggregate1", []Event{e1, e2})

			events := store.GetEventsFrom("aggregate1", e1.Id, 5)

			Expect(events).To(HaveLen(2))
			Expect(events[0].Id).To(Equal(e1.Id))
			Expect(events[1].Id).To(Equal(e2.Id))
		})

	})

	It("Should deserialize payload", func() {
		err := store.SaveEvents("aggregate-1", events())

		Expect(err).To(BeNil())
		Expect(store.GetEvents("aggregate-1")[0].Payload.(EventPayload).Data).To(Equal("payload"))
	})

	It("Should get specific event", func() {
		expected := NewEvent("aid", 0, EventPayload{"payload"})
		err := store.SaveEvents("a-id", []Event{expected})

		actual := store.GetEvent(expected.Id)

		Expect(err).To(BeNil())
		Expect(expected).To(Equal(actual))
	})

})
