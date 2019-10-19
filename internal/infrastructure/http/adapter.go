package http

import (
	"github.com/lzakharov/goss/internal/domain"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

//Adapter represents a HTTP adapter.
type Adapter interface {
	Run() error
	Shutdown() error
}

//NewAdapter creates a new HTTP adapter.
func NewAdapter(logger *zap.Logger, config *Config, service domain.Service) Adapter {
	adapter := &adapter{
		logger:    logger,
		validator: validator.New(),
		config:    config,
		service:   service,
	}

	adapter.server = &fasthttp.Server{
		Handler:     adapter.newRouter(),
		ReadTimeout: config.ReadTimeout,
	}

	return adapter
}

type adapter struct {
	logger    *zap.Logger
	config    *Config
	validator *validator.Validate
	service   domain.Service
	server    *fasthttp.Server
}

//Run starts listening and serving HTTP requests.
func (a *adapter) Run() error {
	a.logger.Info("Starting listening and serving HTTP requests.", zap.String("address", a.config.Address))

	if err := a.server.ListenAndServe(a.config.Address); err != nil {
		a.logger.Error("Error listening and serving HTTP requests!", zap.Error(err))
		return err
	}

	return nil
}

//Shutdown gracefully shuts down the adapter.
func (a *adapter) Shutdown() error {
	if err := a.server.Shutdown(); err != nil {
		a.logger.Error("Error gracefully shutting down the HTTP adapter!", zap.Error(err))
		return err
	}

	return nil
}
