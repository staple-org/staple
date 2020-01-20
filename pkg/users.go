package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/service"
)

// RegisterUser takes a storer and creates a user entry.
func RegisterUser(userHandler service.UserHandlerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the nickname for the token.
		user := &models.User{}
		err := c.Bind(user)
		if user.Email == "" || user.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid username or password",
			})
		}
		if err != nil {
			return err
		}
		if ok, err := userHandler.IsRegistered(*user); ok {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "User already registered.",
			})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return userHandler.Register(*user)
	}
}

// ResetPassword takes a user handler and resets a user's password delievered from the token.
func ResetPassword(userHandler service.UserHandlerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the nickname for the token.
		user := &models.User{}
		err := c.Bind(user)
		if user.Email == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "invalid email",
			})
		}
		if err != nil {
			return err
		}
		if ok, err := userHandler.IsRegistered(*user); !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "User not found",
			})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return userHandler.ResetPassword(*user)
	}
}
