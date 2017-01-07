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
