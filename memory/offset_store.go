package memory

import "github.com/yehohanan7/flux/cqrs"

type InMemoryOffsetStore struct {
	offset int
}

func (store *InMemoryOffsetStore) SaveOffset(value int) error {
	store.offset = value
	return nil
}

func (store *InMemoryOffsetStore) GetLastOffset() (int, error) {
	return store.offset, nil
}

//New in memory offset store
func NewOffsetStore() cqrs.OffsetStore {
	return &InMemoryOffsetStore{offset: 0}
}
