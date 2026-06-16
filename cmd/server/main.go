package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/shreyxnsh/anyx-user-api/config"
	"github.com/shreyxnsh/anyx-user-api/internal/handler"
	"github.com/shreyxnsh/anyx-user-api/internal/logger"
	"github.com/shreyxnsh/anyx-user-api/internal/repository"
	"github.com/shreyxnsh/anyx-user-api/internal/routes"
	"github.com/shreyxnsh/anyx-user-api/internal/service"
)

func main() {

	cfg, err := config.Load()
	if err != nil {

		fmt.Fprintf(os.Stderr, "fatal: failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(cfg.Environment); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: failed to initialise logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting anyx-user-api",
		zap.String("env", cfg.Environment),
		zap.Int("port", cfg.Port),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create connection pool", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("database ping failed — check DATABASE_URL in .env", zap.Error(err))
	}
	logger.Info("database connected successfully")

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if ok := errors.As(err, &e); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	routes.Register(app, userHandler)

	addr := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		logger.Info("server listening", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			logger.Error("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down gracefully…")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("forced shutdown after timeout", zap.Error(err))
	}
	logger.Info("server stopped cleanly")
}
