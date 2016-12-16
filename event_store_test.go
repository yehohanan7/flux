package cqrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveEvents(t *testing.T) {
	store := NewInMemoryEventStore()
	assert.NotNil(t, store)
}
