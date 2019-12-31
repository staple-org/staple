package storage

import "github.com/staple-org/staple/internal/models"

// Storer defines a set of functions for storing staples.
type Storer interface {
	Create(staple models.Staple, userID string) error
	Delete(userID string, stapleID string) error
	Get(userID string, stapleID string) (models.Staple, error)
	List(userID string) ([]models.Staple, error)
}
