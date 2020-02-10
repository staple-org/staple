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
	"github.com/staple-org/staple/pkg/config"
)

// TokenHandler creates a JWT token for a given user.
func TokenHandler(userHandler service.UserHandlerer) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Get the nickname for the token.
		user := &models.User{}
		err := c.Bind(user)
		if err != nil {
			config.Opts.Logger.Error().Err(err).Msg("Failed to bind user")
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
				"message": "username or password mismatch",
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
		t, err := token.SignedString([]byte(config.Opts.GlobalTokenKey))
		if err != nil {
			config.Opts.Logger.Error().Err(err).Msg("Failed to generate token.")
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
			return []byte(config.Opts.GlobalTokenKey), nil
		default:
			return nil, signingMethodError
		}
	})
	if err != nil {
		config.Opts.Logger.Error().Err(err).Msg("Failed to get token")
		return nil, err
	}

	return token, nil
}
