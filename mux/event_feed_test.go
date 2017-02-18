package mux

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/cqrs/cqrs"
	"github.com/yehohanan7/cqrs/feed"
	"github.com/yehohanan7/cqrs/memory"
)

var _ = Describe("Event Feed", func() {
	var (
		router *mux.Router
		store  EventStore
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		store = memory.NewEventStore()
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

var _ = Describe("Atom Feed", func() {

	var (
		router    *mux.Router
		server    *httptest.Server
		store     EventStore
		eventsUrl string
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		store = memory.NewEventStore()
		EventFeed(router, store)
		server = httptest.NewServer(router)
		eventsUrl = fmt.Sprintf("%s/events?format=atom", server.URL)
	})

	It("Publish events as atom feed", func() {
		store.SaveEvents("some_aggregate", []Event{NewEvent("AggregateName", 0, "event payload")})
		request, _ := http.NewRequest("GET", eventsUrl, nil)
		response, err := http.DefaultClient.Do(request)

		Expect(err).Should(BeNil())
		Expect(response).ShouldNot(BeNil())
		Expect(response.StatusCode).To(Equal(http.StatusOK))
		Expect(response.Header.Get("Content-Type")).To(Equal("text/xml"))
	})

})

var _ = Describe("Json Feed", func() {
	var (
		router    *mux.Router
		server    *httptest.Server
		store     EventStore
		eventsUrl string
	)

	BeforeEach(func() {
		router = mux.NewRouter()
		store = memory.NewEventStore()
		EventFeed(router, store)
		server = httptest.NewServer(router)
		eventsUrl = fmt.Sprintf("%s/events?format=json", server.URL)
	})

	It("Should publish events as feeds", func() {
		var feed feed.JsonEventFeed
		event := NewEvent("AggregateName", 0, "event payload")
		store.SaveEvents("some_aggregate", []Event{event})
		request, _ := http.NewRequest("GET", eventsUrl, nil)
		response, err := http.DefaultClient.Do(request)

		Expect(err).Should(BeNil())
		Expect(response).ShouldNot(BeNil())
		Expect(response.StatusCode).To(Equal(http.StatusOK))
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(body, &feed)
		Expect(feed.Description).To(Equal("event feed"))
		Expect(feed.Events).Should(HaveLen(1))
		Expect(feed.Events[0].EventId).To(Equal(event.Id))
		Expect(feed.Events[0].AggregateName).To(Equal("AggregateName"))
		Expect(feed.Events[0].AggregateVersion).To(Equal(0))
		Expect(feed.Events[0].EventType).To(Equal("string"))
		Expect(feed.Events[0].Url).To(Equal(fmt.Sprintf("%s/events/%s", server.URL, event.Id)))
	})
})
