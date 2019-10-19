package security

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/lzakharov/goss/internal/domain"
	"go.uber.org/zap"
)

var (
	timePoint = time.Now()

	config = &Config{
		KeyPrefix:            "auth",
		Secret:               []byte("secret"),
		AccessTokenLifetime:  24 * time.Hour,
		RefreshTokenLifetime: 30 * 24 * time.Hour,
	}

	alice = &domain.User{
		ID:       42,
		Username: "alice",
		Role:     "client",
	}

	aliceAccessTokenClaims = &domain.AccessTokenClaims{
		UserID: alice.ID,
		Role:   alice.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: timePoint.Add(config.AccessTokenLifetime).Unix(),
		},
	}

	aliceRefreshTokenClaims = &domain.RefreshTokenClaims{
		UserID: alice.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: timePoint.Add(config.RefreshTokenLifetime).Unix(),
		},
	}

	aliceAccessToken  string
	aliceRefreshToken string
	aliceAuthData     *domain.AuthData
)

func init() {
	aliceAccessToken, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, aliceAccessTokenClaims).SignedString(config.Secret)
	aliceRefreshToken, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, aliceRefreshTokenClaims).SignedString(config.Secret)
	aliceAuthData = &domain.AuthData{
		AccessToken:  aliceAccessToken,
		ExpiresAt:    timePoint.Add(config.AccessTokenLifetime).Unix(),
		RefreshToken: aliceRefreshToken,
	}
}

func TestNewAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := zap.NewExample()
	redisClient := NewMockRedisClient(ctrl)

	expected := &adapter{
		logger:      logger,
		config:      config,
		redisClient: redisClient,
	}

	actual := NewAdapter(logger, config, redisClient)
	require.Equal(t, expected, actual)
}

func TestAdapter_IsAlive(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := zap.NewExample()

	adapter := &adapter{
		logger: logger,
		config: config,
	}

	t.Run("alive", func(t *testing.T) {
		redisClient := NewMockRedisClient(ctrl)
		redisClient.EXPECT().Ping().Return(redis.NewStatusResult("ok", nil))
		adapter.redisClient = redisClient

		require.True(t, adapter.IsAlive())
	})

	t.Run("dead", func(t *testing.T) {
		redisClient := NewMockRedisClient(ctrl)
		redisClient.EXPECT().Ping().Return(redis.NewStatusResult("", errors.New("redis error")))
		adapter.redisClient = redisClient

		require.False(t, adapter.IsAlive())
	})
}

func TestAdapter_CreateAuthData(t *testing.T) {
	ctrl := gomock.NewController(t)

	patch := monkey.Patch(time.Now, func() time.Time {
		return timePoint
	})
	defer patch.Unpatch()

	logger := zap.NewExample()

	adapter := &adapter{
		logger: logger,
		config: config,
	}

	t.Run("normal", func(t *testing.T) {
		redisClient := NewMockRedisClient(ctrl)
		redisClient.EXPECT().
			Set("auth42", gomock.Any(), config.RefreshTokenLifetime).
			Return(redis.NewStatusResult("ok", nil))

		adapter.redisClient = redisClient

		actual, err := adapter.CreateAuthData(alice)
		require.NoError(t, err)
		require.Equal(t, aliceAuthData, actual)
	})
}

func TestAdapter_GetAccessTokenClaims(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := zap.NewExample()

	adapter := &adapter{
		logger: logger,
		config: config,
	}

	t.Run("normal", func(t *testing.T) {
		authDataJSON, err := json.Marshal(aliceAuthData)
		require.NoError(t, err)

		redisClient := NewMockRedisClient(ctrl)
		redisClient.EXPECT().
			Get("auth42").
			Return(redis.NewStringResult(string(authDataJSON), nil))

		adapter.redisClient = redisClient

		actual, err := adapter.GetAccessTokenClaims(aliceAccessToken)
		require.NoError(t, err)
		require.Equal(t, aliceAccessTokenClaims, actual)
	})
}

func TestAdapter_GetRefreshTokenClaims(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := zap.NewExample()

	adapter := &adapter{
		logger: logger,
		config: config,
	}

	t.Run("normal", func(t *testing.T) {
		authDataJSON, err := json.Marshal(aliceAuthData)
		require.NoError(t, err)

		redisClient := NewMockRedisClient(ctrl)
		redisClient.EXPECT().
			Get("auth42").
			Return(redis.NewStringResult(string(authDataJSON), nil))

		adapter.redisClient = redisClient

		actual, err := adapter.GetRefreshTokenClaims(aliceRefreshToken)
		require.NoError(t, err)
		require.Equal(t, aliceRefreshTokenClaims, actual)
	})
}

func TestAdapter_InvalidateUserAuthData(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := zap.NewExample()

	adapter := &adapter{
		logger: logger,
		config: config,
	}

	t.Run("normal", func(t *testing.T) {
		redisClient := NewMockRedisClient(ctrl)
		redisClient.EXPECT().
			Del("auth42").
			Return(redis.NewIntCmd(1, nil))

		adapter.redisClient = redisClient

		require.NoError(t, adapter.InvalidateUserAuthData(alice.ID))
	})
}
