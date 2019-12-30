package service

import (
	"context"
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
	// TODO: this is probably rather a user ID which is authenticated.
	Create(user *users.User) (staple *Staple, err error)
	Delete(user *users.User, id string) (err error)
	Get(user *users.User, id string) (staple *Staple, err error)
	MarkAsRead(user *users.User, staple *Staple) (err error)
	List(user *users.User) (staples []*Staple, err error)
	Archive(user *users.User, staple *Staple) (err error)
}

// PostgresStapler defines a stapler which stores the staples in Postgres DB.
type PostgresStapler struct {
	// TODO: Add DB connection?
	ctx context.Context
}

// NewPostgresStapler creates a new Postgres based Stapler which will have a connection to a DB.
func NewPostgresStapler() Staplerer {
	return PostgresStapler{ctx: context.Background()}
}

// Create creates a new Staple for the given user.
func (p PostgresStapler) Create(user *users.User) (staple *Staple, err error) {
	panic("implement me")
}

// Delete deletes a given staple for a user.
func (p PostgresStapler) Delete(user *users.User, id string) (err error) {
	panic("implement me")
}

// Get retrieves a Staple for a given user with ID.
func (p PostgresStapler) Get(user *users.User, id string) (staple *Staple, err error) {
	panic("implement me")
}

// MarkAsRead marks a staple as read.
func (p PostgresStapler) MarkAsRead(user *users.User, staple *Staple) (err error) {
	panic("implement me")
}

// List lists all staples for a given user.
func (p PostgresStapler) List(user *users.User) (staples []*Staple, err error) {
	panic("implement me")
}

// Archive will archive a staple which isn't removed but rather not shown in the queue.
// Archived Staples can be retrieved and vewied in any order.
func (p PostgresStapler) Archive(user *users.User, staple *Staple) (err error) {
	panic("implement me")
}
