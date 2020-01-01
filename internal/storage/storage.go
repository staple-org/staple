package storage

import "github.com/staple-org/staple/internal/models"

// Storer defines a set of functions for storing staples.
type Storer interface {
	Create(staple models.Staple, userID string) error
	Delete(user string, stapleID string) error
	Get(user string, stapleID string) (models.Staple, error)
	List(user string) ([]models.Staple, error)
	Archive(user string, stapleID string) error
}
