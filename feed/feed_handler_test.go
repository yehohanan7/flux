package feed

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	"net/http"

	"fmt"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"

	"github.com/yehohanan7/flux/memory"
)

type SamplePayload struct {
	Data string `json:"data"`
}

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
		router.HandleFunc("/events", FeedHandler(store))
		router.HandleFunc("/events/{id}", FeedHandler(store))
		server = httptest.NewServer(router)
		eventsUrl = fmt.Sprintf("%s/events", server.URL)
	})

	It("Should publish event feed", func() {
		var feed JsonEventFeed
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

	It("Should expose event by event id", func() {
		var expected SamplePayload
		event := NewEvent("AggregateName", 0, SamplePayload{"sample data"})
		store.SaveEvents("some_aggregate", []Event{event})

		request, _ := http.NewRequest("GET", eventsUrl+"/"+event.Id, nil)
		response, err := http.DefaultClient.Do(request)

		Expect(err).Should(BeNil())
		Expect(response).ShouldNot(BeNil())
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(body, &expected)
		Expect(expected.Data).To(Equal("sample data"))
	})
})
