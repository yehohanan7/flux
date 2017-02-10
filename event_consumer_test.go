package cqrs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		consumer.EventConsumer = NewConsumer("http://localhost:1212/events", consumer, offsetStore)
	})

	It("Should consume events from another service", func() {
		Expect("").To(Equal(""))
	})
})
