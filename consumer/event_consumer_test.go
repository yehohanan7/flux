package consumer

import (
	"io/ioutil"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
	"github.com/yehohanan7/flux/memory"
	gock "gopkg.in/h2non/gock.v1"
)

type Distance struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type NewStarBorn struct {
	Description       string   `json:"description"`
	DistanceFromEarth Distance `json:"distance_from_earth"`
}

type NewGalaxyFormed struct {
	Description string `json:"description"`
}

type UniverseEventHandler struct {
	EventConsumer
	Stars    int
	Galaxies int
}

func (c *UniverseEventHandler) HandleNewStars(event NewStarBorn) {
	c.Stars++
}

func (c *UniverseEventHandler) HandleGalaxies(event NewGalaxyFormed) {
	c.Galaxies++
}

var _ = Describe("Event Consumer", func() {
	It("Should convert events to event map", func() {
		expected := map[string]reflect.Type{
			"consumer.NewGalaxyFormed": reflect.TypeOf(NewGalaxyFormed{}),
			"consumer.NewStarBorn":     reflect.TypeOf(NewStarBorn{}),
		}

		em := eventMap([]interface{}{NewGalaxyFormed{}, NewStarBorn{}})

		Expect(em).To(Equal(expected))
	})
})

var _ = Describe("Event Consumer", func() {
	baseUrl := "http://localhost:1212"

	BeforeEach(func() {
		feed, _ := ioutil.ReadFile("testdata/universe_events.json")
		star_born, _ := ioutil.ReadFile("testdata/star_born.json")
		galaxy_formed, _ := ioutil.ReadFile("testdata/galaxy_formed.json")

		gock.New(baseUrl).
			Get("/events").
			Reply(200).
			JSON(feed)

		gock.New(baseUrl).
			Get("/events/1").
			Reply(200).
			JSON(string(star_born))

		gock.New(baseUrl).
			Get("/events/2").
			Reply(200).
			JSON(string(galaxy_formed))
	})

	AfterEach(func() {
		gock.Off()
	})

	It("Should consume events", func() {
		events := []interface{}{NewStarBorn{}, NewGalaxyFormed{}}
		consumer := New(baseUrl+"/events", events, memory.NewOffsetStore())
		eventCh, stopCh := make(chan interface{}), make(chan interface{})

		go consumer.Start(eventCh, stopCh)

		Eventually(func() bool {
			d := <-eventCh
			switch d.(type) {
			case NewStarBorn:
				return d.(NewStarBorn).DistanceFromEarth == Distance{5000, "lightyears"} && d.(NewStarBorn).Description == "a new star is born"
			default:
				return false
			}
		}).Should(BeTrue())
	})

})
