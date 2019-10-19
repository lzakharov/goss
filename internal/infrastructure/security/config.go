package security

import "time"

//Config contains a security configuration adapter.
type Config struct {
	KeyPrefix            string             `validate:"required"`
	Secret               []byte             `validate:"required"`
	AccessTokenLifetime  time.Duration      `validate:"required"`
	RefreshTokenLifetime time.Duration      `validate:"required"`
	RedisClient          *RedisClientConfig `validate:"required"`
}
