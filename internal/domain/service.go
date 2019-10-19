package domain

import (
	"time"

	"github.com/lzakharov/goss/internal/version"
	"go.uber.org/zap"
)

//go:generate mockgen -package $GOPACKAGE -source $GOFILE -destination mock_$GOFILE -self_package=github.com/lzakharov/goss/internal/$GOPACKAGE

//Service service represents a service layer.
type Service interface {
	CheckHealth() *Health

	Login(credentials *Credentials) (*AuthData, error)
	RefreshToken(refreshToken string) (*AuthData, error)
	GetUser(userID int64) (*User, error)
	Logout(userID int64) error

	GetAccessTokenClaims(accessToken string) (*AccessTokenClaims, error)
}

//NewService creates a new service.
func NewService(logger *zap.Logger, storage Storage, security Security) Service {
	service := &service{
		logger:   logger,
		storage:  storage,
		security: security,
	}

	return service
}

type service struct {
	logger   *zap.Logger
	storage  Storage
	security Security
}

//CheckHealth checks the application health.
func (s *service) CheckHealth() *Health {
	return &Health{
		Version:  version.Version,
		Storage:  s.storage.IsAlive(),
		Security: s.security.IsAlive(),
	}
}

//Login creates user auth data by the credentials.
func (s *service) Login(credentials *Credentials) (*AuthData, error) {
	user, err := s.storage.GetUserByCredentials(credentials)
	if err != nil {
		s.logger.Error("Error getting user by credentials!",
			zap.String("username", credentials.Username),
			zap.Error(err))
		return nil, err
	}

	authData, err := s.security.CreateAuthData(user)
	if err != nil {
		s.logger.Error("Error creating user auth data!",
			zap.Int64("userID", user.ID),
			zap.Error(err))
		return nil, err
	}

	return authData, nil
}

//RefreshToken refreshes user auth data.
func (s *service) RefreshToken(refreshToken string) (*AuthData, error) {
	claims, err := s.security.GetRefreshTokenClaims(refreshToken)
	if err != nil {
		s.logger.Error("Error refreshing user auth data!",
			zap.String("refreshToken", refreshToken),
			zap.Error(err))
		return nil, err
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.storage.GetUser(claims.UserID)
	if err != nil {
		s.logger.Error("Error getting user by id!",
			zap.Int64("userID", claims.UserID),
			zap.Error(err))
		return nil, err
	}

	authData, err := s.security.CreateAuthData(user)
	if err != nil {
		s.logger.Error("Error creating user auth data!",
			zap.Int64("userID", user.ID),
			zap.Error(err))
		return nil, err
	}

	return authData, nil
}

//GetUser gets a user by id.
func (s *service) GetUser(userID int64) (*User, error) {
	user, err := s.storage.GetUser(userID)
	if err != nil {
		s.logger.Error("Error getting user by id!",
			zap.Int64("userID", userID),
			zap.Error(err))
		return nil, err
	}

	return user, nil
}


//Logout invalidates user's auth data.
func (s *service) Logout(userID int64) error {
	if err := s.security.InvalidateUserAuthData(userID); err != nil {
		s.logger.Error("Error invalidating user's auth data!",
			zap.Int64("userID", userID),
			zap.Error(err))
		return err
	}

	return nil
}

//GetAccessTokenClaims validates the access token and returns claims if success.
func (s *service) GetAccessTokenClaims(accessToken string) (*AccessTokenClaims, error) {
	claims, err := s.security.GetAccessTokenClaims(accessToken)
	if err != nil {
		s.logger.Error("Error getting claims from the access token!",
			zap.String("accessToken", accessToken),
			zap.Error(err))
		return nil, err
	}

	return claims, nil
}
