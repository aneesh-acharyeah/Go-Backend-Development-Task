package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anish/backend-development-task/config"
	"github.com/anish/backend-development-task/db/sqlc"
	"github.com/anish/backend-development-task/internal/handler"
	"github.com/anish/backend-development-task/internal/logger"
	appmiddleware "github.com/anish/backend-development-task/internal/middleware"
	"github.com/anish/backend-development-task/internal/models"
	"github.com/anish/backend-development-task/internal/repository"
	"github.com/anish/backend-development-task/internal/routes"
	"github.com/anish/backend-development-task/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	appLogger, err := logger.New(cfg.Environment)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = appLogger.Sync()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("failed to create database pool", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		appLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	appLogger.Info("database connection established")

	app := fiber.New(fiber.Config{
		AppName:      "users-api",
		ErrorHandler: fiberErrorHandler(appLogger),
	})

	app.Use(appmiddleware.RequestLogger(appLogger))

	queries := sqlc.New(pool)
	userRepository := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, appLogger)
	routes.Register(app, userHandler)

	go shutdownOnSignal(app, appLogger)

	address := ":" + cfg.Port
	appLogger.Info("server starting", zap.String("address", address))
	if err := app.Listen(address); err != nil {
		appLogger.Fatal("server stopped unexpectedly", zap.Error(err))
	}
}

func fiberErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		message := "internal server error"
		code := "internal_error"

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			status = fiberErr.Code
			message = fiberErr.Message
			code = errorCodeForStatus(status)
		} else {
			logger.Error("unhandled application error", zap.Error(err))
		}

		return c.Status(status).JSON(models.ErrorResponse{
			Error:   code,
			Message: message,
		})
	}
}

func errorCodeForStatus(status int) string {
	switch {
	case status == fiber.StatusNotFound:
		return "not_found"
	case status >= fiber.StatusBadRequest && status < fiber.StatusInternalServerError:
		return "bad_request"
	default:
		return "internal_error"
	}
}

func shutdownOnSignal(app *fiber.App, logger *zap.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received")
	if err := app.Shutdown(); err != nil {
		logger.Error("server shutdown failed", zap.Error(err))
	}
}
