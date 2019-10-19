package domain

import "errors"

var (
	//ErrInternalStorage represents the internal storage error.
	ErrInternalStorage = errors.New("internal storage error")

	//ErrNotFound represents the object not found error.
	ErrNotFound = errors.New("not found")

	//ErrInvalidCredentials represents the invalid credentials error.
	ErrInvalidCredentials  = errors.New("invalid credentials")

	//ErrInternalSecurity represents the internal security error.
	ErrInternalSecurity    = errors.New("internal security error")

	//ErrInvalidAccessToken represents the invalid access token error.
	ErrInvalidAccessToken  = errors.New("invalid access token")

	//ErrInvalidRefreshToken represents the invalid refresh token error.
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
