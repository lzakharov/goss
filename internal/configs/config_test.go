package configs

import (
	"os"
	"testing"
	"time"

	"github.com/lzakharov/goss/internal/infrastructure/security"

	"github.com/stretchr/testify/require"
	"github.com/lzakharov/goss/internal/infrastructure/http"
	"github.com/lzakharov/goss/internal/infrastructure/storage"
	"go.uber.org/zap"
)

func TestNewConfig(t *testing.T) {
	logger := zap.NewExample()

	t.Run("with a valid environment", func(t *testing.T) {
		require.NoError(t, os.Setenv("APP_STORAGE_DB_DRIVERNAME", "postgres"))
		require.NoError(t, os.Setenv("APP_STORAGE_DB_DSN", "postgres"))
		require.NoError(t, os.Setenv("APP_STORAGE_DB_MAXOPENCONNS", "64"))
		require.NoError(t, os.Setenv("APP_STORAGE_DB_MAXIDLECONNS", "64"))
		require.NoError(t, os.Setenv("APP_STORAGE_DB_CONNMAXLIFETIME", "30s"))
		require.NoError(t, os.Setenv("APP_STORAGE_DB_MIGRATIONS_DIALECT", "postgres"))
		require.NoError(t, os.Setenv("APP_STORAGE_DB_MIGRATIONS_DIR", "./migrations/dev/postgres"))

		require.NoError(t, os.Setenv("APP_SECURITY_KEYPREFIX", "auth"))
		require.NoError(t, os.Setenv("APP_SECURITY_SECRET", "secret"))
		require.NoError(t, os.Setenv("APP_SECURITY_ACCESSTOKENLIFETIME", "2h"))
		require.NoError(t, os.Setenv("APP_SECURITY_REFRESHTOKENLIFETIME", "720h"))
		require.NoError(t, os.Setenv("APP_SECURITY_REDISCLIENT_ADDR", "redis:6379"))

		require.NoError(t, os.Setenv("APP_HTTP_ADDRESS", "127.0.0.1:8080"))
		require.NoError(t, os.Setenv("APP_HTTP_READTIMEOUT", "5s"))
		expected := &Config{
			Storage: &storage.Config{
				DB: &storage.DBConfig{
					DriverName:      "postgres",
					DSN:             "postgres",
					MaxOpenConns:    64,
					MaxIdleConns:    64,
					ConnMaxLifeTime: 30 * time.Second,
					Migrations: &storage.MigrationsConfig{
						Dialect: "postgres",
						Dir:     "./migrations/dev/postgres",
					},
				},
			},
			Security: &security.Config{
				KeyPrefix:            "auth",
				Secret:               []byte("secret"),
				AccessTokenLifetime:  2 * time.Hour,
				RefreshTokenLifetime: 30 * 24 * time.Hour,
				RedisClient: &security.RedisClientConfig{
					Addr: "redis:6379",
				},
			},
			HTTP: &http.Config{
				Address:     "127.0.0.1:8080",
				ReadTimeout: 5 * time.Second,
			},
		}

		actual, err := NewConfig(logger)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("with an empty environment", func(t *testing.T) {
		require.NoError(t, os.Setenv("APP_HTTP_ADDRESS", ""))
		require.NoError(t, os.Setenv("APP_HTTP_READTIMEOUT", ""))
		_, err := NewConfig(logger)
		require.Error(t, err)
	})
}
