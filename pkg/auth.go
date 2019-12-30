package pkg

import "github.com/labstack/echo/v4"

// AuthZeroMiddleware defines an authentication middleware using Auth 0.
func AuthZeroMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// c.Error(errors.New("unauthorized") // if auth0 returns 0... ?
		return next(c)
	}
}

// AuthZeroCallback handles the callback from Auth0 service.
func AuthZeroCallback() echo.HandlerFunc {
	return func(context echo.Context) error {
		// Hande auth zero callback.
		return nil
	}
}
