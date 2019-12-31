package pkg

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"

	"github.com/staple-org/staple/internal/service"
	"github.com/staple-org/staple/pkg/auth"
)

// Serve starts the Stapler API server.
func Serve() error {
	log.Println("Starting listener...")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	sessionSecret := os.Getenv("STAPLE_SESSION_SECRET")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"https://staple.cronohub.org"},
	//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	//}))

	api := "/rest/api/1"
	gob.Register(map[string]interface{}{})

	// Landing page
	e.GET("/login", auth.LoginHandler())
	e.POST("/logout", func(context echo.Context) error {
		return nil
	})
	e.GET("/callback", auth.Callback())

	stapler := service.NewPostgresStapler()
	g := e.Group(api+"/staple", auth.Middleware)
	g.POST("/", AddStaple(stapler))
	g.POST("/:id/archive", AddStaple(stapler))
	g.POST("/:id/markasread", AddStaple(stapler))
	g.GET("/:id", AddStaple(stapler))
	g.DELETE("/:id", DeleteStaple(stapler))
	g.GET("/", AddStaple(stapler))

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
