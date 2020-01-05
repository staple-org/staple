package service

import "github.com/staple-org/staple/internal/models"

// UserHandlerer defines a service which can manage users.
type UserHandlerer interface {
	Register(user models.User) error
	Delete(user models.User) error
	ResetPassword(user models.User) error
	IsRegistered(user models.User) (ok bool, err error)
	GetToken(user models.User) (token string, err error)
}
