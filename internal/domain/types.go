package domain

import "github.com/dgrijalva/jwt-go"

//Health contains an application health status.
type Health struct {
	Version  string `json:"version"`
	Storage  bool   `json:"storage"`
	Security bool   `json:"security"`
}

//Credentials contains user credentials.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//AuthData contains auth data.
type AuthData struct {
	AccessToken  string `json:"accessToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	RefreshToken string `json:"refreshToken"`
}

//AccessTokenClaims contains access token claims.
type AccessTokenClaims struct {
	UserID int64  `json:"userID"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

//RefreshTokenClaims contains refresh token claims.
type RefreshTokenClaims struct {
	UserID int64 `json:"userID"`
	jwt.StandardClaims
}

//User contains a user.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
