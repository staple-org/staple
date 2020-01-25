package storage

import (
	"github.com/staple-org/staple/internal/models"
)

// PostgresUserStorer is a storer which uses Postgres as a storage backend.
type InMemoryUserStorer struct {
	Err   error
	store map[string]*models.User
}

// NewInMemoryUserStorer creates a new Postgres storage medium.
func NewInMemoryUserStorer() InMemoryUserStorer {
	return InMemoryUserStorer{
		store: make(map[string]*models.User),
	}
}

// Create saves a user in the db.
func (s InMemoryUserStorer) Create(email string, password []byte) error {
	s.store[email] = &models.User{
		Email:       email,
		Password:    string(password),
		ConfirmCode: "",
		MaxStaples:  DefaultMaxStaples,
	}
	return s.Err
}

// Delete deletes a user from the db.
func (s InMemoryUserStorer) Delete(email string) error {
	if s.Err != nil {
		return s.Err
	}
	delete(s.store, email)
	return s.Err
}

// Get retrieves a user.
func (s InMemoryUserStorer) Get(email string) (*models.User, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	return s.store[email], s.Err
}

// Update updates a user with a given email address.
func (s InMemoryUserStorer) Update(email string, newUser models.User) error {
	if s.Err != nil {
		return s.Err
	}
	s.store[email] = &newUser
	return s.Err
}
