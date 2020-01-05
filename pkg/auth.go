package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/staple-org/staple/internal/service"

	"github.com/staple-org/staple/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// TokenHandler creates a JWT token for a given user.
func TokenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {

		// TODO: Verify that this user exists.
		// Get the nickname for the token.
		user := &models.User{}
		err := c.Bind(user)
		if err != nil {
			return err
		}
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = user.Email // from context
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(Opts.GlobalTokenKey))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}

// GetToken gets the JWT token from the echo context
func GetToken(c echo.Context) (*jwt.Token, error) {
	// Get the token
	jwtRaw := c.Request().Header.Get("Authorization")
	split := strings.Split(jwtRaw, " ")
	if len(split) != 2 {
		return nil, errors.New("unauthorized")
	}
	jwtString := split[1]
	// Parse token
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		signingMethodError := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		switch token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return []byte(Opts.GlobalTokenKey), nil
		default:
			return nil, signingMethodError
		}
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// RegisterUser takes a storer and creates a user entry.
func RegisterUser(userHandler service.UserHandlerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := models.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}
		if ok, err := userHandler.IsRegistered(u); ok {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "User already registered.",
			})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return nil
	}
}
