package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/VikaGo/REST_API/config"
	"github.com/VikaGo/REST_API/controller"
	"github.com/VikaGo/REST_API/logger"
	Error "github.com/VikaGo/REST_API/pkg/error"
	"github.com/VikaGo/REST_API/pkg/validator"
	"github.com/VikaGo/REST_API/service"
	"github.com/VikaGo/REST_API/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// config
	cfg := config.Get()

	// logger
	l := logger.Get()

	// Init repository store
	store, err := store.New(ctx)
	if err != nil {
		return errors.Wrap(err, "store.New failed")
	}

	// Init service manager
	serviceManager, err := service.NewManager(ctx, store)
	if err != nil {
		return errors.Wrap(err, "manager.New failed")
	}

	// Init controllers
	userController := controller.NewUsers(ctx, serviceManager, l)

	// Initialize Echo instance
	e := echo.New()
	e.Validator = validator.NewValidator()
	e.HTTPErrorHandler = Error.Error

	// Disable Echo JSON logger in debug mode
	if cfg.LogLevel == "debug" {
		if l, ok := e.Logger.(*echoLog.Logger); ok {
			l.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line}")
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API V1
	v1 := e.Group("/v1")

	// User routes
	userRoutes := v1.Group("/users")
	userRoutes.POST("/login", userController.LogIn)
	userRoutes.GET("/:id", userController.Get)
	userRoutes.DELETE("/:id", userController.Delete)
	userRoutes.PUT("/:id", userController.Update)

	// Start server
	s := &http.Server{
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
		Addr:         ":8080",
	}
	e.Logger.Fatal(e.StartServer(s))

	return nil
}
