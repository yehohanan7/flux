package cqrs

import "testing"
import "github.com/stretchr/testify/assert"

type TestProjectionEvent struct {
	Data string
}

func (e TestProjectionEvent) Trigger() string {
	return "unknown"
}

type TestProjection struct {
	Projection
	Data string
}

func (p *TestProjection) Handle(event TestProjectionEvent) {
	p.Data = event.Data
}

func TestBuildProjection(t *testing.T) {
	model := new(TestProjection)
	model.Projection = NewProjection(model)

	model.Apply([]interface{}{TestProjectionEvent{"some data"}})

	assert.Equal(t, model.Data, "some data")
}
