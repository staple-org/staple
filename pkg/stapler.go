package pkg

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/service"
)

// AddStaple creates a staple using a stapler and a given user.
func AddStaple(stapler service.Staplerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("auth-session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Session not found.")
		}
		profile := sess.Values["profile"].(map[string]interface{})
		user := &models.User{
			Email: profile["nickname"].(string),
		}
		// TODO: Construct staple here. POST will have the information needed.
		err = stapler.Create(models.Staple{}, user)
		if err != nil {
			apiError := ApiError("Unable to create staple for user.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		return c.NoContent(http.StatusOK)
	}
}

// GetNext retrieves the oldest entry from the list which is not archived.
func GetNext(staple service.Staplerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

// GetStaple retrieves a single staple based on an ID.
func GetStaple(stapler service.Staplerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := GetToken(c)
		if err != nil {
			return err
		}
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		userModel := &models.User{
			Email: email,
		}
		id := c.Param("id")
		if id == "" {
			return errors.New("invalid id")
		}
		s, err := stapler.Get(userModel, id)
		if err != nil {
			apiError := ApiError("something went wrong", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		if s == nil {
			apiError := ApiError("staple not found", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		var staple = struct {
			Staple models.Staple `json:"staple"`
		}{
			Staple: *s,
		}
		return c.JSON(http.StatusOK, staple)
	}
}

// ListStaples will list all staples which belong to a user.
func ListStaples(stapler service.Staplerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := GetToken(c)
		if err != nil {
			return err
		}
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		userModel := &models.User{
			Email: email,
		}
		s, err := stapler.List(userModel)
		if err != nil {
			apiError := ApiError("Unable to list staples for user.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		var staples = struct {
			Staples []models.Staple `json:"staples"`
		}{
			Staples: s,
		}
		return c.JSON(http.StatusOK, staples)
	}
}

// DeleteStaple deteles a staple with a given ID.
func DeleteStaple(stapler service.Staplerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get user ID from context.. Call delete.
		return stapler.Delete(nil, "")
	}
}
