package memory

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
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
	})

	It("Should reject events which has wrong aggregate version", func() {
		err1 := store.SaveEvents("aggregate-1", events())
		err2 := store.SaveEvents("aggregate-1", events())
		err3 := store.SaveEvents("aggregate-1", []Event{NewEvent("sample_aggregate", 3, EventPayload{"payload"})})

		Expect(err1).To(BeNil())
		Expect(err2).ShouldNot(BeNil())
		Expect(err3).To(BeNil())
	})

	var _ = Describe("Fetching all event metadata from a secific offset", func() {
		It("Should get the events", func() {
			e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
			e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})
			e3 := NewEvent("sample_aggregate", 3, EventPayload{"payload"})
			e4 := NewEvent("sample_aggregate", 4, EventPayload{"payload"})

			store.SaveEvents("aggregate1", []Event{e1, e2, e3, e4})

			metas := store.GetEventMetaDataFrom(1, 2)
			Expect(metas).To(HaveLen(2))
			Expect(metas[0].Id).To(Equal(e2.Id))
			Expect(metas[1].Id).To(Equal(e3.Id))
		})

		It("Should handle count gracefully", func() {
			e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
			e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})

			store.SaveEvents("aggregate1", []Event{e1, e2})

			metas := store.GetEventMetaDataFrom(1, 5)

			Expect(metas).To(HaveLen(1))
			Expect(metas[0].Id).To(Equal(e2.Id))
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
