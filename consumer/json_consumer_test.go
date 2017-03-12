package consumer

import (
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
	"github.com/yehohanan7/flux/memory"
	. "github.com/yehohanan7/flux/utils"
	"gopkg.in/h2non/gock.v1"
)

var _ = Describe("Event Consumer", func() {
	var (
		handler     *UniverseEventHandler
		offsetStore OffsetStore
	)

	BeforeEach(func() {
		offsetStore = memory.NewOffsetStore()
		handler = new(UniverseEventHandler)
		handler.EventConsumer = NewEventConsumer("http://localhost:1212/events", handler, offsetStore)
	})

	AfterEach(func() {
		gock.Off()
	})

	It("Should consume events from another service", func() {
		feed, _ := ioutil.ReadFile("testdata/universe_events.json")
		star_born, _ := ioutil.ReadFile("testdata/star_born.json")
		galaxy_formed, _ := ioutil.ReadFile("testdata/galaxy_formed.json")
		gock.New("http://localhost:1212").
			Get("/events").
			Reply(200).
			JSON(string(feed))

		gock.New("http://localhost:1212").
			Get("/events/1").
			Reply(200).
			JSON(string(star_born))

		gock.New("http://localhost:1212").
			Get("/events/2").
			Reply(200).
			JSON(string(galaxy_formed))

		handler.Start()

		WaitUntil(func() bool { return handler.Stars == 1 }, 20*time.Second)
		Expect(handler.Stars).To(Equal(1))
	})

})
