package cqrs

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

type NewStarBorn struct {
}

type NewGalaxyFormed struct {
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
		json, _ := ioutil.ReadFile("testdata/universe_events.json")
		gock.New("http://localhost:1212").
			Get("/events").
			Reply(200).
			JSON(string(json))

		resp, _ := http.Get("http://localhost:1212/events")
		body, _ := ioutil.ReadAll(resp.Body)
		Expect(string(body)).To(Equal(""))
	})

})
