package http

import (
	"github.com/lzakharov/goss/internal/domain"
)

var errStatus = map[error]int{
	domain.ErrInvalidCredentials:  4011,
	domain.ErrInvalidAccessToken:  4012,
	domain.ErrInvalidRefreshToken: 4013,
	domain.ErrNotFound:            4040,

	domain.ErrInternalStorage:  5001,
	domain.ErrInternalSecurity: 5002,
}

//ErrorResponse is an error response.
type ErrorResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"requestID"`
}
