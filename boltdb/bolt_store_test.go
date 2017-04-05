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
		e1 := NewEvent("sample_aggregate", 1, EventPayload{"payload"})
		e2 := NewEvent("sample_aggregate", 2, EventPayload{"payload"})

		err := store.SaveEvents("aggregate-1", []Event{e1, e2})

		Expect(err).To(BeNil())
		Expect(store.GetEvent(e1.Id).Id).To(Equal(e1.Id))
		Expect(store.GetEvent(e2.Id).Id).To(Equal(e2.Id))
	})
})
