package service

import (
	"time"

	"github.com/staple-org/staple/internal/users"
)

// Staple defines a Staple in the system.
type Staple struct {
	Name            string    `json:"name"`
	ID              string    `json:"id"`
	Content         []byte    `json:"content"`
	CreateTimestamp time.Time `json:"create_timestamp"`
	Archived        bool      `json:"archived"`
}

// Staplerer should be able to do that following:
// Create Staple for a user
// Delete Staple for a user
// Get a Staple for a user
// Archive Staple
// MarkAsRead
// List Staples for a user
type Staplerer interface {
	Create(user *users.User) (staple *Staple, err error)
	Delete(user *users.User, staple *Staple) (ok bool, err error)
	Get(user *users.User) (staple *Staple, err error)
	MarkAsRead(user *users.User, staple *Staple) (err error)
	List(user *users.User) (staples []*Staple, err error)
	Archive(user *users.User, staple *Staple) (err error)
}
