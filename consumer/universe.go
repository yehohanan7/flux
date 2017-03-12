package consumer

import (
	. "github.com/yehohanan7/flux/cqrs"
)

type NewStarBorn struct {
	Description string `json:"description"`
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
