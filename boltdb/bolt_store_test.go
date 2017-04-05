package boltdb

import (
	"encoding/gob"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
)

const DB_PATH = "flux_test.db"

type EventPayload struct {
	Data string
}

var _ = Describe("Bolt Event Store", func() {

	var store EventStore

	BeforeEach(func() {
		gob.Register(EventPayload{})
		store = NewBoltStore(DB_PATH)
	})

	AfterEach(func() {
		os.Remove(DB_PATH)
	})

	It("Should save events", func() {
		expected := NewEvent("sample_aggregate", 1, EventPayload{"payload"})

		err := store.SaveEvents("aggregate-1", []Event{expected})

		actual := store.GetEvent(expected.Id)
		Expect(err).To(BeNil())
		Expect(actual.Id).To(Equal(expected.Id))
		Expect(actual.Payload).To(Equal(expected.Payload))
	})
})
