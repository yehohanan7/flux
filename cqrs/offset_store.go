package cqrs

type OffsetStore interface {
	SaveOffset(int) error
	GetLastOffset() (int, error)
}
