package pkg

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/staple-org/staple/internal/service"
)

// AddStaple creates a staple using a stapler and a given user.
func AddStaple(stapler service.Staplerer) echo.HandlerFunc {
	return func(context echo.Context) error {
		s, err := stapler.Create(nil)
		log.Printf("%+v", s)
		return err
	}
}

// DeleteStaple deteles a staple with a given ID.
func DeleteStaple(stapler service.Staplerer) echo.HandlerFunc {
	return func(context echo.Context) error {
		// Get user ID from context.. Call delete.
		return stapler.Delete(nil, "")
	}
}
