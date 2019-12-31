package service

import (
	"context"
	"fmt"

	"github.com/staple-org/staple/internal/models"

	"github.com/staple-org/staple/internal/storage"
)

// Staplerer should be able to do that following:
// Create Staple for a user
// Delete Staple for a user
// Get a Staple for a user
// Archive Staple
// MarkAsRead
// List Staples for a user
type Staplerer interface {
	Create(staple models.Staple, user *models.User) (err error)
	Delete(user *models.User, id string) (err error)
	Get(user *models.User, id string) (staple models.Staple, err error)
	MarkAsRead(user *models.User, staple models.Staple) (err error)
	List(user *models.User) (staples []models.Staple, err error)
	Archive(user *models.User, staple models.Staple) (err error)
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
func (p Stapler) Create(staple models.Staple, user *models.User) error {
	fmt.Println("Staple created for user: ", user)
	return nil
}

// Delete deletes a given staple for a user.
func (p Stapler) Delete(user *models.User, id string) (err error) {
	fmt.Println("Staple Delete called.")
	return nil
}

// Get retrieves a Staple for a given user with ID.
func (p Stapler) Get(user *models.User, id string) (models.Staple, error) {
	fmt.Println("Staple Get called.")
	return models.Staple{}, nil
}

// MarkAsRead marks a staple as read.
func (p Stapler) MarkAsRead(user *models.User, staple models.Staple) error {
	fmt.Println("Staple Mark as Read called.")
	return nil
}

// List lists all staples for a given user.
func (p Stapler) List(user *models.User) ([]models.Staple, error) {
	fmt.Println("Staple List called.")
	list, err := p.storer.List(user.ID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Archive will archive a staple which isn't removed but rather not shown in the queue.
// Archived Staples can be retrieved and vewied in any order.
func (p Stapler) Archive(user *models.User, staple models.Staple) error {
	fmt.Println("Staple Archive called.")
	return nil
}
