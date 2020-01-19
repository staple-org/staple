package storage

import "github.com/staple-org/staple/internal/models"

// StapleStorer defines a set of functions for storing staples.
type StapleStorer interface {
	Create(staple models.Staple, email string) error
	Delete(email string, stapleID int) error
	Get(email string, stapleID int) (*models.Staple, error)
	List(email string) ([]models.Staple, error)
	Archive(email string, stapleID int) error
	Oldest(email string) (*models.Staple, error)
	ShowArchive(email string) ([]models.Staple, error)
}

// UserStorer defines a set of functions for storing users.
type UserStorer interface {
	Create(email string, password []byte) error
	Delete(email string) error
	Get(email string) (*models.User, error)
	Update(email string, newUser models.User) error
}
