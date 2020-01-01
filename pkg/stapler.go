package pkg

import (
	"net/http"

	"github.com/staple-org/staple/internal/models"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

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
			Nickname: profile["nickname"].(string),
			ID:       profile["aud"].(string),
		}
		// TODO: Construct staple here. POST will have the information needed.
		err = stapler.Create(models.Staple{}, user)
		if err != nil {
			var message = struct {
				code    int
				message string
			}{
				code:    http.StatusInternalServerError,
				message: "Unable to create staple for user.",
			}
			return c.JSON(http.StatusInternalServerError, message)
		}
		return c.NoContent(http.StatusOK)
	}
}

// ListStaples will list all staples which belong to a user.
func ListStaples(stapler service.Staplerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("auth-session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Session not found.")
		}
		profile := sess.Values["profile"].(map[string]interface{})
		user := &models.User{
			Nickname: profile["nickname"].(string),
			ID:       profile["aud"].(string),
		}
		s, err := stapler.List(user)
		if err != nil {
			var message = struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, message)
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
