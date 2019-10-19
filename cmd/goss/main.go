package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lzakharov/goss/internal/configs"
	"github.com/lzakharov/goss/internal/domain"
	"github.com/lzakharov/goss/internal/infrastructure/http"
	"github.com/lzakharov/goss/internal/infrastructure/security"
	"github.com/lzakharov/goss/internal/infrastructure/storage"
	"github.com/lzakharov/goss/internal/version"
	"go.uber.org/zap"
)

const (
	projectEnv  = "ENVIRONMENT"
	development = "DEV"
	production  = "PROD"
)

func newLogger() (*zap.Logger, error) {
	env := os.Getenv(projectEnv)

	switch env {
	case development:
		return zap.NewDevelopment()
	case production:
		return zap.NewProduction()
	default:
		return nil, fmt.Errorf("unknown environment '%s', check the '%s' environment variable", env, projectEnv)
	}
}

func main() {
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}
	logger = logger.Named(version.Version)

	config, err := configs.NewConfig(logger)
	if err != nil {
		logger.Panic("Error creating a new configuration!", zap.Error(err))
	}
	logger.Debug("Configuration read successfully.", zap.Any("config", config))

	db, err := storage.NewDB(config.Storage.DB)
	if err != nil {
		logger.Panic("Error creating a new SQL database!", zap.Error(err))
	}

	storageAdapter := storage.NewAdapter(logger, db)

	redisClient, err := security.NewRedisClient(config.Security.RedisClient)
	if err != nil {
		logger.Panic("Error creating a new Redis client!", zap.Error(err))
	}

	securityAdapter := security.NewAdapter(logger, config.Security, redisClient)

	service := domain.NewService(logger, storageAdapter, securityAdapter)

	httpAdapter := http.NewAdapter(logger, config.HTTP, service)

	shutdown := make(chan error, 1)

	go func(shutdown chan<- error) {
		if err := httpAdapter.Run(); err != nil {
			shutdown <- err
		}
	}(shutdown)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		logger.Info("Got the signal!", zap.Any("signal", x))
	case err := <-shutdown:
		logger.Error("Error running the application!", zap.Error(err))
	}

	logger.Info("Stopping the application...")

	if err := httpAdapter.Shutdown(); err != nil {
		logger.Panic("Error shutting down the HTTP adapter!", zap.Error(err))
	}

	logger.Info("The application gracefully stopped.")
}
