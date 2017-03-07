package consumer

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/feed"
	"github.com/yehohanan7/flux/memory"
	gock "gopkg.in/h2non/gock.v1"
)

type SampleEvent struct {
}

var _ = Describe("Event Consumer", func() {
	baseUrl := "http://localhost:1212"

	AfterEach(func() {
		gock.Off()
	})

	It("Should consume events", func() {
		feed, _ := ioutil.ReadFile("testdata/universe_events.json")
		gock.New(baseUrl).
			Get("/events").
			Reply(200).
			JSON(feed)

		consumer := New(baseUrl+"/events", []interface{}{}, memory.NewOffsetStore())
		eventCh, stopCh := make(chan interface{}), make(chan interface{})

		go consumer.Start(eventCh, stopCh)

		Eventually(func() string {
			d := <-eventCh
			feed := d.(JsonEventFeed)
			return feed.Description
		}).Should(Equal("event feed"))
	})
})
