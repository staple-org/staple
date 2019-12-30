package pkg

import (
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

// Serve starts the Stapler API server.
func Serve() error {
	log.Println("Starting listener...")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := "/rest/api/1"

	e.POST(api+"/staple", func(context echo.Context) error {
		return nil
	})


	hostPort := fmt.Sprintf("%s:%s", Opts.Hostname, Opts.Port)

	// Start TLS with certificate paths
	if len(Opts.ServerKeyPath) > 0 && len(Opts.ServerCrtPath) > 0 {
		e.Pre(middleware.HTTPSRedirect())
		return e.StartTLS(hostPort, Opts.ServerCrtPath, Opts.ServerKeyPath)
	}

	// Start Auto TLS server
	if Opts.AutoTLS {
		if len(Opts.CacheDir) < 1 {
			return errors.New("cache dir must be provided if autoTLS is enabled")
		}
		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.Cache = autocert.DirCache(Opts.CacheDir)
		return e.StartAutoTLS(hostPort)
	}
	// Start regular server
	return e.Start(hostPort)
}
