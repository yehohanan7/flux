package cqrs

import "testing"
import "github.com/stretchr/testify/assert"

func TestEventKey(t *testing.T) {
	key := buildKey("aggregateId", NewEvent("AggregateName", 1, "payload"))

	assert.Equal(t, "aggregateId_1", key)
}
