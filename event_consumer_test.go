package cqrs

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

type NewStarBorn struct {
	Description string `json:"description"`
}

type NewGalaxyFormed struct {
	Description string `json:"description"`
}

type UniverseEventConsumer struct {
	EventConsumer
	Stars    int
	Galaxies int
}

func (c *UniverseEventConsumer) HandleNewStars(event NewStarBorn) {
	c.Stars++
}

func (c *UniverseEventConsumer) HandleGalaxies(event NewGalaxyFormed) {
	c.Galaxies++
}

var _ = Describe("Event Consumer", func() {
	var (
		consumer    *UniverseEventConsumer
		offsetStore OffsetStore
	)

	BeforeEach(func() {
		offsetStore = NewInMemoryOffsetStore()
		consumer = new(UniverseEventConsumer)
		consumer.EventConsumer = NewEventConsumer("http://localhost:1212/events", consumer, offsetStore)
	})

	AfterEach(func() {
		gock.Off()
	})

	It("Should consume events from another service", func() {
		feed, _ := ioutil.ReadFile("testdata/universe_events.json")
		event, _ := ioutil.ReadFile("testdata/star_born.json")
		gock.New("http://localhost:1212").
			Get("/events").
			Reply(200).
			JSON(string(feed))

		gock.New("http://localhost:1212").
			Get("/events/1").
			Reply(200).
			JSON(string(event))

		gock.New("http://localhost:1212").
			Get("/events/2").
			Reply(200).
			JSON(string(event))

		gock.New("http://localhost:1212").
			Get("/events/3").
			Reply(200).
			JSON(string(event))

		consumer.Start()

		Expect(consumer.Stars).To(Equal(2))
	})

})
