package service

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/google/uuid"
	"github.com/staple-org/staple/pkg/config"
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
	SendConfirmCode(user models.User) error
	VerifyConfirmCode(user models.User) (bool, error)
	SetMaximumStaples(user models.User, maxStaples int) error
	GetMaximumStaples(user models.User) (int, error)
	ChangePassword(user models.User, newPassword string) error
}

// UserHandler defines a storage using user handler.
type UserHandler struct {
	ctx      context.Context
	store    storage.UserStorer
	notifier Notifier
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
// This happens after the confirmation was successfull.
func (u UserHandler) ResetPassword(user models.User) error {
	// get the stored user based on the provided email.
	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		return err
	}
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
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
	storedUser.Password = string(hashPassword)
	storedUser.ConfirmCode = ""
	if err := u.store.Update(storedUser.Email, *storedUser); err != nil {
		return err
	}

	return u.notifier.Notify(storedUser.Email, PasswordReset, newPassword)
}

// SendConfirmCode sends a confirm code which has to be verified.
func (u UserHandler) SendConfirmCode(user models.User) error {
	confirmUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		return err
	}

	storedUser.ConfirmCode = confirmUUID.String()
	if err := u.store.Update(storedUser.Email, *storedUser); err != nil {
		return err
	}
	return u.notifier.Notify(storedUser.Email, GenerateConfirmCode, storedUser.ConfirmCode)
}

// VerifyConfirmCode will match the confirm code with a stored code for an email address.
// If the match is successful the code is removed and the password is reset.
func (u UserHandler) VerifyConfirmCode(user models.User) (ok bool, err error) {
	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		return false, err
	}
	if storedUser == nil {
		return false, errors.New("user not found")
	}
	if user.ConfirmCode == storedUser.ConfirmCode && user.Email == storedUser.Email {
		if err := u.ResetPassword(user); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("confirm code did not match")
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

// SetMaximumStaples sets the user's maximum number of allowed staples.
func (u UserHandler) SetMaximumStaples(user models.User, maxStaples int) error {
	if maxStaples <= 0 || maxStaples > 100 {
		return errors.New("invalid staple setting")
	}
	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		return err
	}
	storedUser.MaxStaples = maxStaples
	if err := u.store.Update(user.Email, *storedUser); err != nil {
		return err
	}
	return nil
}

// ChangePassword changes the user's password to a new given string.
func (u UserHandler) ChangePassword(user models.User, newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}
	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		config.Opts.Logger.Error().Err(err).Msg("Error while getting user")
		return err
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	storedUser.Password = string(hashPassword)
	if err := u.store.Update(user.Email, *storedUser); err != nil {
		config.Opts.Logger.Error().Err(err).Msg("Error while storing user")
		return err
	}
	return nil
}

// GetMaximumStaples returns the maximum allowed configured staples for a user.
func (u UserHandler) GetMaximumStaples(user models.User) (staples int, err error) {
	storedUser, err := u.store.Get(user.Email)
	if err != nil {
		return 0, err
	}
	if storedUser == nil {
		return 0, nil
	}
	return storedUser.MaxStaples, nil
}

// NewUserHandler creates a new user handler.
func NewUserHandler(ctx context.Context, store storage.UserStorer, notifier Notifier) UserHandler {
	return UserHandler{
		ctx:      ctx,
		store:    store,
		notifier: notifier,
	}
}
