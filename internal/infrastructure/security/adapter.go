package security

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/lzakharov/goss/internal/domain"
	"go.uber.org/zap"
)

//NewAdapter creates a new security adapter.
func NewAdapter(logger *zap.Logger, config *Config, redisClient RedisClient) domain.Security {
	adapter := &adapter{
		logger:      logger,
		config:      config,
		redisClient: redisClient,
	}

	return adapter
}

type adapter struct {
	logger      *zap.Logger
	config      *Config
	redisClient RedisClient
}

//IsAlive returns true if the adapter can ping redis.
func (a *adapter) IsAlive() bool {
	return a.redisClient.Ping().Err() == nil
}

//CreateAuthData generates auth data for the specified user.
func (a *adapter) CreateAuthData(user *domain.User) (*domain.AuthData, error) {
	now := time.Now()

	accessToken, err := a.newAccessToken(now, user)
	if err != nil {
		a.logger.Error("Error creating a new access token!",
			zap.Int64("userID", user.ID),
			zap.Error(err))
		return nil, domain.ErrInternalSecurity
	}

	refreshToken, err := a.newRefreshToken(now, user.ID)
	if err != nil {
		a.logger.Error("Error creating a new refresh token!",
			zap.Int64("userID", user.ID),
			zap.Error(err))
		return nil, domain.ErrInternalSecurity
	}

	authData := &domain.AuthData{
		AccessToken:  accessToken,
		ExpiresAt:    now.Add(a.config.AccessTokenLifetime).Unix(),
		RefreshToken: refreshToken,
	}

	key := a.newKey(user.ID)
	value, err := json.Marshal(authData)
	if err != nil {
		a.logger.Error("Error saving user's auth data!",
			zap.Int64("userID", user.ID),
			zap.Error(err))
		return nil, domain.ErrInternalSecurity
	}

	if err := a.redisClient.Set(key, value, a.config.RefreshTokenLifetime).Err(); err != nil {
		a.logger.Error("Error saving user's auth data!",
			zap.Int64("userID", user.ID),
			zap.Error(err))
		return nil, domain.ErrInternalSecurity
	}

	return authData, nil
}

//GetAccessTokenClaims gets access token claims.
func (a *adapter) GetAccessTokenClaims(accessToken string) (*domain.AccessTokenClaims, error) {
	claims := new(domain.AccessTokenClaims)

	token, err := jwt.ParseWithClaims(accessToken, claims, a.jwtKeyFunc)
	if err != nil {
		a.logger.Error("Error parsing access token!",
			zap.String("accessToken", accessToken),
			zap.Error(err))
		return nil, domain.ErrInvalidAccessToken
	}
	if !token.Valid {
		a.logger.Warn("Invalid access token!",
			zap.String("accessToken", accessToken),
			zap.Error(err))
		return nil, domain.ErrInvalidAccessToken
	}

	key := a.newKey(claims.UserID)
	data, err := a.redisClient.Get(key).Result()
	if err != nil {
		a.logger.Error("Error getting access token!",
			zap.String("key", key),
			zap.Error(err))

		if err == redis.Nil {
			return nil, domain.ErrInvalidAccessToken
		}
		return nil, domain.ErrInternalSecurity
	}

	authData := new(domain.AuthData)

	if err := json.Unmarshal([]byte(data), &authData); err != nil {
		a.logger.Error("Error getting access token!",
			zap.String("key", key),
			zap.Error(err))
		return nil, domain.ErrInternalSecurity
	}

	if accessToken != authData.AccessToken {
		a.logger.Warn("Transferred token does not match stored!", zap.String("key", key))
		return nil, domain.ErrInvalidAccessToken
	}

	return claims, nil
}

//GetAccessTokenClaims gets refresh token claims.
func (a *adapter) GetRefreshTokenClaims(refreshToken string) (*domain.RefreshTokenClaims, error) {
	claims := new(domain.RefreshTokenClaims)

	token, err := jwt.ParseWithClaims(refreshToken, claims, a.jwtKeyFunc)
	if err != nil {
		a.logger.Error("Error parsing refresh token!",
			zap.String("refreshToken", refreshToken),
			zap.Error(err))
		return nil, domain.ErrInvalidRefreshToken
	}
	if !token.Valid {
		a.logger.Warn("Invalid refresh token!",
			zap.String("refreshToken", refreshToken),
			zap.Error(err))
		return nil, domain.ErrInvalidRefreshToken
	}

	key := a.newKey(claims.UserID)
	data, err := a.redisClient.Get(key).Result()
	if err != nil {
		a.logger.Error("Error getting refresh token!",
			zap.String("key", key),
			zap.Error(err))

		if err == redis.Nil {
			return nil, domain.ErrInvalidRefreshToken
		}
		return nil, domain.ErrInternalSecurity
	}

	authData := new(domain.AuthData)

	if err := json.Unmarshal([]byte(data), &authData); err != nil {
		a.logger.Error("Error getting refresh token!",
			zap.String("key", key),
			zap.Error(err))
		return nil, domain.ErrInternalSecurity
	}

	if refreshToken != authData.RefreshToken {
		a.logger.Warn("Transferred token does not match stored!", zap.String("key", key))
		return nil, domain.ErrInvalidRefreshToken
	}

	return claims, nil
}

//InvalidateUserAuthData invalidates user's auth data.
func (a *adapter) InvalidateUserAuthData(userID int64) error {
	key := a.newKey(userID)
	if _, err := a.redisClient.Del(key).Result(); err != nil {
		a.logger.Error("Error deleting user's auth data!",
			zap.String("key", key),
			zap.Error(err))
		return domain.ErrInternalSecurity
	}

	return nil
}

func (a *adapter) newAccessToken(now time.Time, user *domain.User) (string, error) {
	claims := &domain.AccessTokenClaims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(a.config.AccessTokenLifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(a.config.Secret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a *adapter) newRefreshToken(now time.Time, userID int64) (string, error) {
	claims := &domain.RefreshTokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(a.config.RefreshTokenLifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString(a.config.Secret)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (a *adapter) newKey(userID int64) string {
	return fmt.Sprintf(keyFormat, a.config.KeyPrefix, userID)
}

func (a *adapter) jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return a.config.Secret, nil
}
