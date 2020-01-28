package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Register(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	err := userHandler.Register(u)
	assert.NoError(t, err)

	// verifiy by getting the created user
	ok, err := userHandler.IsRegistered(u)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestUserHandler_ChangePassword(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	err := userHandler.Register(u)
	assert.NoError(t, err)

	// verifiy by getting the created user
	ok, err := userHandler.IsRegistered(u)
	assert.NoError(t, err)
	assert.True(t, ok)

	// changing password
	err = userHandler.ChangePassword(u, "newPassword")
	assert.NoError(t, err)

	// verify password match
	ok, err = userHandler.PasswordMatch(u)
	assert.Error(t, err)
	assert.False(t, ok)

	// verify password match with new
	u.Password = "newPassword"
	ok, err = userHandler.PasswordMatch(u)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestUserHandler_ChangePassword_NoNewPassword(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}

	// changing password
	err := userHandler.ChangePassword(u, "")
	assert.EqualError(t, err, "password cannot be empty")
}

func TestUserHandler_Delete(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	err := userHandler.Register(u)
	assert.NoError(t, err)

	// verifiy by getting the created user
	ok, err := userHandler.IsRegistered(u)
	assert.NoError(t, err)
	assert.True(t, ok)

	// Delete the user
	err = userHandler.Delete(u)
	assert.NoError(t, err)

	// verifiy by getting the created user
	ok, err = userHandler.IsRegistered(u)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestUserHandler_PasswordMatch(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}
	err := userHandler.Register(u)
	assert.NoError(t, err)

	ok, err := userHandler.PasswordMatch(u)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestUserHandler_PasswordMatch_UserNotFound(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "",
		MaxStaples:  25,
	}

	_, err := userHandler.PasswordMatch(u)
	assert.EqualError(t, err, "user not found")
}

func TestUserHandler_ResetPassword(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "11111",
		MaxStaples:  25,
	}
	err := userHandler.Register(u)
	assert.NoError(t, err)

	err = userHandler.ResetPassword(u)
	assert.NoError(t, err)

	// Verify that the password no longer match
	ok, err := userHandler.PasswordMatch(u)
	assert.Error(t, err)
	assert.False(t, ok)

	// Verify that the notifier sent out the new password
	body := notifier.buffer.String()
	assert.NotEmpty(t, body)

	var newPassword string
	_, _ = fmt.Sscanf(body, `Dear test@test.com
Your password has been successfully reset to: %s. Please change as soon as possible.`, &newPassword)
	u.Password = newPassword[0 : len(newPassword)-1] // trim the "." from the end of the string.

	// Verify that the new password matches
	ok, err = userHandler.PasswordMatch(u)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestUserHandler_SetMaximumStaples(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "11111",
		MaxStaples:  25,
	}
	err := userHandler.Register(u)
	assert.NoError(t, err)
	err = userHandler.SetMaximumStaples(u, 10)
	assert.NoError(t, err)
	n, err := userHandler.GetMaximumStaples(u)
	assert.NoError(t, err)
	assert.Equal(t, 10, n)
}

func TestUserHandler_SendConfirmCode(t *testing.T) {
	store := storage.NewInMemoryUserStorer()
	notifier := NewBufferNotifier()
	userHandler := NewUserHandler(context.Background(), store, notifier)

	u := models.User{
		Email:       "test@test.com",
		Password:    "password",
		ConfirmCode: "11111",
		MaxStaples:  25,
	}
	err := userHandler.Register(u)
	assert.NoError(t, err)
	err = userHandler.SendConfirmCode(u)
	assert.NoError(t, err)

	body := notifier.buffer.String()
	var code string
	_, _ = fmt.Sscanf(body, `Dear test@test.com
Please enter the following code into the confirm code window: %s`, &code)
	u.ConfirmCode = code
	ok, err := userHandler.VerifyConfirmCode(u)
	assert.NoError(t, err)
	assert.True(t, ok)
}
