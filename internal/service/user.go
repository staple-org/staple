package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/storage"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

// UserHandlerer defines a service which can manage users.
type UserHandlerer interface {
	Register(user models.User) error
	Delete(user models.User) error
	ResetPassword(user models.User) error
	IsRegistered(user models.User) (ok bool, err error)
	PasswordMatch(user models.User) (ok bool, err error)
}

// UserHandler defines a storage using user handler.
type UserHandler struct {
	ctx   context.Context
	store storage.UserStorer
}

// Register registers a user.
func (u UserHandler) Register(user models.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = u.store.Create(user.Email, hashPassword)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a user.
func (u UserHandler) Delete(user models.User) error {
	if ok, _ := u.IsRegistered(user); !ok {
		return errors.New("user not found")
	}
	if ok, err := u.PasswordMatch(user); !ok {
		return errors.New("password did not match")
	} else if err != nil {
		return err
	}
	return u.store.Delete(user.Email)
}

// ResetPassword generates a new password for a user and send it via email.
func (u UserHandler) ResetPassword(user models.User) error {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	newPassword := string(bytes)
	hashPassword, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fmt.Println("Actually set newPassword:", hashPassword)

	return SendResetPasswordEmail(user.Email, newPassword)
}

// IsRegistered checks if a user exists in the system.
func (u UserHandler) IsRegistered(user models.User) (ok bool, err error) {
	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		return false, err
	}
	if storedUser == nil {
		return false, nil
	}
	return true, nil
}

// PasswordMatch checks if a stored password matches that of a given one.
func (u UserHandler) PasswordMatch(user models.User) (ok bool, err error) {
	plain := []byte(user.Password)

	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		return false, err
	}
	if storedUser == nil {
		return false, errors.New("user not found")
	}

	hash := []byte(storedUser.Password)
	err = bcrypt.CompareHashAndPassword(hash, plain)
	if err != nil {
		return false, err
	}
	return true, nil
}

// NewUserHandler creates a new user handler.
func NewUserHandler(ctx context.Context, store storage.UserStorer) UserHandler {
	return UserHandler{
		ctx:   ctx,
		store: store,
	}
}
