package domain

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/golang/mock/gomock"
	"github.com/lzakharov/goss/internal/version"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := zap.NewExample()
	storage := NewMockStorage(ctrl)
	security := NewMockSecurity(ctrl)
	expected := &service{
		logger:   logger,
		storage:  storage,
		security: security,
	}

	actual := NewService(logger, storage, security)
	require.Equal(t, expected, actual)
}

func TestService_CheckHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := zap.NewExample()
	storage := NewMockStorage(ctrl)
	storage.EXPECT().IsAlive().Return(true)
	security := NewMockSecurity(ctrl)
	security.EXPECT().IsAlive().Return(true)
	service := &service{
		logger:   logger,
		storage:  storage,
		security: security,
	}
	expected := &Health{
		Version:  version.Version,
		Storage:  true,
		Security: true,
	}

	actual := service.CheckHealth()
	require.Equal(t, expected, actual)
}

func TestService_Login(t *testing.T) {
	t.Run("with valid credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		credentials := &Credentials{
			Username: "user",
			Password: "password",
		}

		user := &User{
			ID:       1,
			Username: "user",
			Role:     "client",
		}

		expected := &AuthData{
			AccessToken:  "accessToken",
			ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
			RefreshToken: "refreshToken",
		}

		logger := zap.NewExample()
		storage := NewMockStorage(ctrl)
		storage.EXPECT().GetUserByCredentials(credentials).Return(user, nil)
		security := NewMockSecurity(ctrl)
		security.EXPECT().CreateAuthData(user).Return(expected, nil)
		service := &service{
			logger:   logger,
			storage:  storage,
			security: security,
		}

		actual, err := service.Login(credentials)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("with invalid credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		credentials := &Credentials{
			Username: "user",
			Password: "password",
		}

		expected := ErrInvalidCredentials

		logger := zap.NewExample()
		storage := NewMockStorage(ctrl)
		storage.EXPECT().GetUserByCredentials(credentials).Return(nil, expected)
		service := &service{
			logger:  logger,
			storage: storage,
		}

		_, err := service.Login(credentials)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})

	t.Run("with broken storage", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		credentials := &Credentials{
			Username: "user",
			Password: "password",
		}

		expected := ErrInternalStorage

		logger := zap.NewExample()
		storage := NewMockStorage(ctrl)
		storage.EXPECT().GetUserByCredentials(credentials).Return(nil, expected)
		service := &service{
			logger:  logger,
			storage: storage,
		}

		_, err := service.Login(credentials)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})

	t.Run("with broken security", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		credentials := &Credentials{
			Username: "user",
			Password: "password",
		}

		user := &User{
			ID:       42,
			Username: "user",
			Role:     "client",
		}

		expected := ErrInternalSecurity

		logger := zap.NewExample()
		storage := NewMockStorage(ctrl)
		storage.EXPECT().GetUserByCredentials(credentials).Return(user, nil)
		security := NewMockSecurity(ctrl)
		security.EXPECT().CreateAuthData(user).Return(nil, expected)
		service := &service{
			logger:   logger,
			storage:  storage,
			security: security,
		}

		_, err := service.Login(credentials)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})
}

func TestService_RefreshToken(t *testing.T) {
	t.Run("with valid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		refreshToken := "refreshToken"

		claims := &RefreshTokenClaims{
			UserID: 42,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
			},
		}

		user := &User{
			ID:       42,
			Username: "alice",
			Role:     "client",
		}

		expected := &AuthData{
			AccessToken:  "newAccessToken",
			ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
			RefreshToken: "newRefreshToken",
		}

		logger := zap.NewExample()
		storage := NewMockStorage(ctrl)
		storage.EXPECT().GetUser(int64(42)).Return(user, nil)
		security := NewMockSecurity(ctrl)
		security.EXPECT().GetRefreshTokenClaims(refreshToken).Return(claims, nil)
		security.EXPECT().CreateAuthData(user).Return(expected, nil)
		service := &service{
			logger:   logger,
			security: security,
			storage:  storage,
		}

		actual, err := service.RefreshToken(refreshToken)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("with expired token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		refreshToken := "refreshToken"

		claims := &RefreshTokenClaims{
			UserID: 42,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(-time.Hour).Unix(),
			},
		}

		expected := ErrInvalidRefreshToken

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().GetRefreshTokenClaims(refreshToken).Return(claims, nil)
		service := &service{
			logger:   logger,
			security: security,
		}

		_, err := service.RefreshToken(refreshToken)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})

	t.Run("with invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		refreshToken := "refreshToken"

		expected := ErrInvalidRefreshToken

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().GetRefreshTokenClaims(refreshToken).Return(nil, expected)
		service := &service{
			logger:   logger,
			security: security,
		}

		_, err := service.RefreshToken(refreshToken)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})

	t.Run("with broken storage", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		refreshToken := "refreshToken"

		claims := &RefreshTokenClaims{
			UserID: 42,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
			},
		}

		expected := ErrInternalStorage

		logger := zap.NewExample()
		storage := NewMockStorage(ctrl)
		storage.EXPECT().GetUser(int64(42)).Return(nil, expected)
		security := NewMockSecurity(ctrl)
		security.EXPECT().GetRefreshTokenClaims(refreshToken).Return(claims, nil)
		service := &service{
			logger:   logger,
			security: security,
			storage:  storage,
		}

		_, err := service.RefreshToken(refreshToken)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})

	t.Run("with broken security", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		refreshToken := "refreshToken"

		expected := ErrInternalSecurity

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().GetRefreshTokenClaims(refreshToken).Return(nil, expected)
		service := &service{
			logger:   logger,
			security: security,
		}

		_, err := service.RefreshToken(refreshToken)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})
}

func TestService_GetUser(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		expected := &User{
			ID:       42,
			Username: "username",
			Role:     "user",
		}

		logger := zap.NewExample()
		storage := NewMockStorage(ctrl)
		storage.EXPECT().GetUser(expected.ID).Return(expected, nil)
		service := &service{
			logger:  logger,
			storage: storage,
		}

		actual, err := service.GetUser(42)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestService_Logout(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userID := int64(42)

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().InvalidateUserAuthData(userID).Return(nil)
		service := &service{
			logger:   logger,
			security: security,
		}

		require.NoError(t, service.Logout(userID))
	})

	t.Run("with broken security", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userID := int64(42)
		expected := ErrInternalSecurity

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().InvalidateUserAuthData(userID).Return(expected)
		service := &service{
			logger:   logger,
			security: security,
		}

		err := service.Logout(userID)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})
}

func TestService_GetClaims(t *testing.T) {
	t.Run("with valid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userID := int64(42)

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().InvalidateUserAuthData(userID).Return(nil)
		service := &service{
			logger:   logger,
			security: security,
		}

		require.NoError(t, service.Logout(userID))
	})

	t.Run("with invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userID := int64(42)
		expected := ErrInvalidAccessToken

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().InvalidateUserAuthData(userID).Return(expected)
		service := &service{
			logger:   logger,
			security: security,
		}

		err := service.Logout(userID)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})

	t.Run("with broken security", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userID := int64(42)
		expected := ErrInternalSecurity

		logger := zap.NewExample()
		security := NewMockSecurity(ctrl)
		security.EXPECT().InvalidateUserAuthData(userID).Return(expected)
		service := &service{
			logger:   logger,
			security: security,
		}

		err := service.Logout(userID)
		require.Error(t, err)
		require.Equal(t, expected, err)
	})
}
