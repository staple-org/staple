package pkg

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/staple-org/staple/internal/models"
	"github.com/staple-org/staple/internal/service"
)

// AddStaple creates a staple using a stapler and a given user.
// The following properties are enough:
// name, content
func AddStaple(stapler service.Staplerer) echo.HandlerFunc {
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
		staple := &models.Staple{}
		err = c.Bind(staple)
		if err != nil {
			apiError := ApiError("failed to bind body", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		staple.CreatedAt = time.Now().UTC()
		// ID needs to be sequential.
		err = stapler.Create(*staple, userModel)
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
		token, err := GetToken(c)
		if err != nil {
			return err
		}
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		userModel := &models.User{
			Email: email,
		}
		s, err := staple.GetNext(userModel)
		if err != nil {
			apiError := ApiError("failed getting next staple", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		var staple = struct {
			Staple *models.Staple `json:"staple"`
		}{
			Staple: s,
		}
		return c.JSON(http.StatusOK, staple)
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
		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := ApiError("failed to convert id to number", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		s, err := stapler.Get(userModel, n)
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

// ShowArchive returns the archived staples of a user.
func ShowArchive(stapler service.Staplerer) echo.HandlerFunc {
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
		s, err := stapler.ShowArchive(userModel)
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
		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := ApiError("failed to convert id to number", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		err = stapler.Delete(userModel, n)
		if err != nil {
			apiError := ApiError("Unable to delete staple.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		return c.NoContent(http.StatusOK)
	}
}

// ArchiveStaple archives a staple with a given ID.
func ArchiveStaple(stapler service.Staplerer) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get user ID from context.. Call delete.
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
		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := ApiError("failed to convert id to number", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		err = stapler.Archive(userModel, n)
		if err != nil {
			apiError := ApiError("Unable to delete staple.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		return c.NoContent(http.StatusOK)
	}
}
