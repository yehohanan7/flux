package memory

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

func NewOffsetStore() *InMemoryOffsetStore {
	return &InMemoryOffsetStore{offset: -1}
}
