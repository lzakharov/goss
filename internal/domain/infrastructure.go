package domain

//go:generate mockgen -package $GOPACKAGE -source $GOFILE -destination mock_$GOFILE -self_package=github.com/lzakharov/goss/internal/domain/$GOPACKAGE

//Mortal represents a service that can be up or down.
type Mortal interface {
	IsAlive() bool
}

//Storage represents a storage adapter.
type Storage interface {
	Mortal

	GetUser(userID int64) (*User, error)
	GetUserByCredentials(credentials *Credentials) (*User, error)
}

//Security represents a security adapter.
type Security interface {
	Mortal

	CreateAuthData(user *User) (*AuthData, error)
	GetAccessTokenClaims(accessToken string) (*AccessTokenClaims, error)
	GetRefreshTokenClaims(refreshToken string) (*RefreshTokenClaims, error)
	InvalidateUserAuthData(userID int64) error
}
