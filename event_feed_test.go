package cqrs

import (
	"net/http/httptest"

	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Feed validations", func() {
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

var _ = Describe("Event Feed", func() {

	var (
		router    *mux.Router
		server    *httptest.Server
		store     EventStore
		eventsUrl string
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		store = NewEventStore()
		EventFeed(router, store)
		server = httptest.NewServer(router)
		eventsUrl = fmt.Sprintf("%s/events", server.URL)
	})

	It("Publish events as atom feed", func() {
		store.SaveEvents("some_aggregate", []Event{NewEvent("AggregateName", 0, "event payload")})
		request, _ := http.NewRequest("GET", eventsUrl, nil)
		response, err := http.DefaultClient.Do(request)
		//defer response.Body.Close()

		Expect(err).Should(BeNil())
		Expect(response).ShouldNot(BeNil())
		Expect(response.StatusCode).To(Equal(http.StatusOK))
		//body, _ := ioutil.ReadAll(response.Body)
		//Expect(string(body)).To(Equal(""))
	})

})
