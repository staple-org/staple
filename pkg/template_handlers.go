package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Index renders and shows the index.tmpl file.
func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.tmpl", "")
}
