package cqrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	aggregateId := "aggregate-id"
	repo := NewAggregateRepo()
	aggregate := NewAggregate(aggregateId, new(TestEntity), repo)

	err := repo.Save(aggregate)

	assert.Nil(t, err)
	assert.Equal(t, aggregateId, repo.Get(aggregateId).Id)
}
