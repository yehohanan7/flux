package cqrs

import (
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Feed", func() {
	var (
		router *mux.Router
		store  EventStore
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		store = NewEventStore()
	})

	It("Should not accept invalid page size", func() {
		Expect(func() { EventFeed(router, store, -1) }).Should(Panic())
	})

	It("Should use default page size", func() {
		EventFeed(router, store)
		Expect(pageSize).To(Equal(DEFAULT_PAGE_SIZE))
	})

	It("Should configure page size", func() {
		EventFeed(router, store, 5)
		Expect(pageSize).To(Equal(5))
	})

})
