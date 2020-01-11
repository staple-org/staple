package pkg

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"

	"github.com/staple-org/staple/internal/service"
	"github.com/staple-org/staple/internal/storage"
)

// Serve starts the Stapler API server.
func Serve() error {
	log.Println("Starting listener...")
	// Echo instance
	e := echo.New()
	// Register the template renderer

	// Setup Global Token Key
	if Opts.GlobalTokenKey == "" {
		log.Print("Please set a global secret key... Randomly generating one for now...")
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		state := base64.StdEncoding.EncodeToString(b)
		Opts.GlobalTokenKey = state
		log.Println("done.")
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"https://staple.cronohub.org"},
	//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	//}))

	// Register a user.
	ctx := context.Background()
	postgresUserStorer := storage.NewPostgresUserStorer()
	userHandler := service.NewUserHandler(ctx, postgresUserStorer)
	e.POST("/register", RegisterUser(userHandler))
	// Generate a token for a given username.
	e.POST("/get-token", TokenHandler(userHandler))

	api := "/rest/api/1"
	//gob.Register(map[string]interface{}{})
	postgresStapleStorer := storage.NewPostgresStapleStorer()
	stapler := service.NewStapler(postgresStapleStorer)

	// REST api group
	g := e.Group(api+"/staple", middleware.JWT([]byte(Opts.GlobalTokenKey)))
	g.POST("", AddStaple(stapler))
	g.POST("/:id/archive", AddStaple(stapler))
	g.POST("/:id/markasread", AddStaple(stapler))
	g.GET("/:id", GetStaple(stapler))
	g.DELETE("/:id", DeleteStaple(stapler))
	g.GET("", ListStaples(stapler))

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
