package security

import (
	"time"

	"github.com/go-redis/redis"
)

//go:generate mockgen -package $GOPACKAGE -source $GOFILE -destination mock_$GOFILE -self_package=github.com/lzakharov/goss/internal/infrastructure/$GOPACKAGE

//RedisClient represents a redis client.
type RedisClient interface {
	Ping() *redis.StatusCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Del(keys ...string) *redis.IntCmd
}

// RedisClientConfig contains a redis factory configuration.
type RedisClientConfig struct {
	Addr     string `validate:"required"`
	Password string
	DB       int
}

// NewRedisClient creates a new redis client.
func NewRedisClient(config *RedisClientConfig) (RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
