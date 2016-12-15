package cqrs

type AggregateRepository interface {
	Save(aggregate Aggregate) error
	Get(id string) Aggregate
}

type InMemoryAggregateRepo struct {
	aggregates map[string]Aggregate
}

func (repo *InMemoryAggregateRepo) Get(id string) Aggregate {
	return repo.aggregates[id]
}

func (repo *InMemoryAggregateRepo) Save(aggregate Aggregate) error {
	repo.aggregates[aggregate.Id] = aggregate
	return nil
}

func NewInMemoryRepository() AggregateRepository {
	return &InMemoryAggregateRepo{make(map[string]Aggregate)}
}

func NewAggregateRepo() AggregateRepository {
	return NewInMemoryRepository()
}
