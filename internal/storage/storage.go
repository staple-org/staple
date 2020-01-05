package storage

import "github.com/staple-org/staple/internal/models"

// StapleStorer defines a set of functions for storing staples.
type StapleStorer interface {
	Create(staple models.Staple, email string) error
	Delete(email string, stapleID string) error
	Get(email string, stapleID string) (models.Staple, error)
	List(email string) ([]models.Staple, error)
	Archive(email string, stapleID string) error
}

// UserStorer defines a set of functions for storing users.
type UserStorer interface {
	Create(email string, password []byte) error
	Delete(email string) error
	Get(email string) (*models.User, error)
	PasswordMatches(email string, password []byte) (bool, error)
	Update(email string, newUser models.User) error
}
