package configs

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/lzakharov/goss/internal/infrastructure/http"
	"github.com/lzakharov/goss/internal/infrastructure/security"
	"github.com/lzakharov/goss/internal/infrastructure/storage"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const prefix = "app"

//Config contains an application configuration.
type Config struct {
	Storage  *storage.Config  `validate:"required"`
	Security *security.Config `validate:"required"`
	HTTP     *http.Config     `validate:"required"`
}

// NewConfig reads an application configuration from the environment variables.
func NewConfig(logger *zap.Logger) (*Config, error) {
	config := new(Config)

	if err := envconfig.Process(prefix, config); err != nil {
		logger.Error("Error reading an application configuration from the environment variables!", zap.Error(err))
		return nil, err
	}

	if err := validator.New().Struct(config); err != nil {
		logger.Error("Error validating an application configuration!", zap.Error(err))
		return nil, err
	}
	return config, nil
}
