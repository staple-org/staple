package pkg

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"

	"github.com/staple-org/staple/internal/service"
	"github.com/staple-org/staple/internal/storage"
	"github.com/staple-org/staple/pkg/config"
)

// Serve starts the Stapler API server.
func Serve() error {
	// Echo instance
	e := echo.New()
	// Register the template renderer

	// Setup Global Token Key
	if config.Opts.GlobalTokenKey == "" {
		config.Opts.Logger.Info().Msg("Please set a global secret key... Randomly generating one for now...")
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		state := base64.StdEncoding.EncodeToString(b)
		config.Opts.GlobalTokenKey = state
	}

	// Setup Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Opts.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	config.Opts.Logger = zerolog.New(os.Stdout)

	if e := config.Opts.Logger.Debug(); e.Enabled() {
		config.Opts.Logger.Debug().Interface("config", config.Opts)
	}

	config.Opts.Logger.Info().Msg("Starting listener...")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup front-end if not in production mode.
	if !config.Opts.DevMode {
		staticAssets, err := rice.FindBox("../frontend/build")
		if err != nil {
			log.Fatal("Cannot find assets in production")
			return err
		}

		// Register handler for static assets
		assetHandler := http.FileServer(staticAssets.HTTPBox())
		e.GET("/", echo.WrapHandler(assetHandler))
		e.GET("/favicon.ico", echo.WrapHandler(assetHandler))
		e.GET("/static/css/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/static/js/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/static/media/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	}
	// Register a user.
	ctx := context.Background()
	postgresUserStorer := storage.NewPostgresUserStorer()
	emailNotifier := service.NewEmailNotifier()
	userHandler := service.NewUserHandler(ctx, postgresUserStorer, emailNotifier)
	api := "/rest/api/1"

	e.POST(api+"/register", RegisterUser(userHandler))
	// Generate a token for a given username.
	e.POST(api+"/get-token", TokenHandler(userHandler))

	// Reset Password Flow
	e.POST(api+"/reset", ResetPassword(userHandler))
	e.POST(api+"/verify", VerfiyConfirmCode(userHandler))

	//gob.Register(map[string]interface{}{})
	postgresStapleStorer := storage.NewPostgresStapleStorer()
	stapler := service.NewStapler(postgresStapleStorer)

	// REST api group
	g := e.Group(api+"/staple", middleware.JWT([]byte(config.Opts.GlobalTokenKey)))
	g.POST("", AddStaple(stapler, userHandler))
	g.POST("/:id/archive", ArchiveStaple(stapler))
	g.GET("/:id", GetStaple(stapler))
	g.GET("/next", GetNext(stapler))
	g.DELETE("/:id", DeleteStaple(stapler))
	g.GET("/archive", ShowArchive(stapler))
	g.GET("", ListStaples(stapler))

	u := e.Group(api+"/user", middleware.JWT([]byte(config.Opts.GlobalTokenKey)))
	u.POST("/change-password", ChangePassword(userHandler))
	u.POST("/max-staples", SetMaximumStaples(userHandler))
	u.GET("/max-staples", GetMaximumStaples(userHandler))

	hostPort := fmt.Sprintf("%s:%s", config.Opts.Hostname, config.Opts.Port)
	// Start TLS with certificate paths
	if len(config.Opts.ServerKeyPath) > 0 && len(config.Opts.ServerCrtPath) > 0 {
		e.Pre(middleware.HTTPSRedirect())
		return e.StartTLS(hostPort, config.Opts.ServerCrtPath, config.Opts.ServerKeyPath)
	}

	// Start Auto TLS server
	if config.Opts.AutoTLS {
		if len(config.Opts.CacheDir) < 1 {
			return errors.New("cache dir must be provided if autoTLS is enabled")
		}
		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.Cache = autocert.DirCache(config.Opts.CacheDir)
		return e.StartAutoTLS(hostPort)
	}
	// Start regular server
	return e.Start(hostPort)
}
