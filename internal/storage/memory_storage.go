package storage

import "github.com/staple-org/staple/internal/models"

// InMemoryStorer is a storage device used for testing purposes.
type InMemoryStorer struct {
	// email as key
	stapleStore map[string][]models.Staple
	// users
	userStore []models.User
}

// NewInMemoryStorer returns a new in memory storage service
func NewInMemoryStorer() InMemoryStorer {
	return InMemoryStorer{
		stapleStore: make(map[string][]models.Staple),
		userStore:   make([]models.User, 0),
	}
}
