package service

import (
	"context"
	"fmt"
	"time"

	"github.com/staple-org/staple/internal/storage"
	"github.com/staple-org/staple/internal/users"
)

// Staple defines a Staple in the system.
type Staple struct {
	Name            string    `json:"name"`
	ID              string    `json:"id"`
	Content         string    `json:"content"`
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
	Create(user *users.User) (staple Staple, err error)
	Delete(user *users.User, id string) (err error)
	Get(user *users.User, id string) (staple Staple, err error)
	MarkAsRead(user *users.User, staple Staple) (err error)
	List(user *users.User) (staples []Staple, err error)
	Archive(user *users.User, staple Staple) (err error)
}

// Stapler defines a stapler which stores the staples in Postgres DB.
type Stapler struct {
	ctx    context.Context
	storer storage.Storer
}

// NewStapler creates a new Postgres based Stapler which will have a connection to a DB.
func NewStapler(storer storage.Storer) Stapler {
	return Stapler{ctx: context.Background(), storer: storer}
}

// Create creates a new Staple for the given user.
func (p Stapler) Create(user *users.User) (staple Staple, err error) {
	fmt.Println("Staple created for user: ", user)
	return Staple{}, nil
}

// Delete deletes a given staple for a user.
func (p Stapler) Delete(user *users.User, id string) (err error) {
	fmt.Println("Staple Delete called.")
	return nil
}

// Get retrieves a Staple for a given user with ID.
func (p Stapler) Get(user *users.User, id string) (staple Staple, err error) {
	fmt.Println("Staple Get called.")
	return Staple{}, nil
}

// MarkAsRead marks a staple as read.
func (p Stapler) MarkAsRead(user *users.User, staple Staple) (err error) {
	fmt.Println("Staple Mark as Read called.")
	return nil
}

// List lists all staples for a given user.
func (p Stapler) List(user *users.User) (staples []Staple, err error) {
	fmt.Println("Staple List called.")
	list, err := p.storer.List(user.ID)
	if err != nil {
		return nil, err
	}
	fmt.Println(list)
	prop := []Staple{
		{
			Name:            "Sample",
			ID:              "ASDF",
			Content:         "This is the content.",
			CreateTimestamp: time.Now(),
			Archived:        false,
		},
	}
	return prop, nil
}

// Archive will archive a staple which isn't removed but rather not shown in the queue.
// Archived Staples can be retrieved and vewied in any order.
func (p Stapler) Archive(user *users.User, staple Staple) (err error) {
	fmt.Println("Staple Archive called.")
	return nil
}
