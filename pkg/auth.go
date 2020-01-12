package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/service"
)

// TokenHandler creates a JWT token for a given user.
func TokenHandler(userHandler service.UserHandlerer) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Get the nickname for the token.
		user := &models.User{}
		err := c.Bind(user)
		if err != nil {
			return err
		}
		if user.Email == "" || user.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid username or password",
			})
		}

		if ok, _ := userHandler.IsRegistered(*user); !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "user not found",
			})
		}

		if ok, err := userHandler.PasswordMatch(*user); !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "password mismatch",
			})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "error while getting password: " + err.Error(),
			})
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
