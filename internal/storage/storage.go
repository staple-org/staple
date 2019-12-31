package storage

// Storer defines a set of functions for storing staples.
type Storer interface {
	Create(userID string) error
	Delete(userID string, stapleID string) error
	Get(userID string, stapleID string) ([]byte, error)
	List(userID string) ([]byte, error)
}
