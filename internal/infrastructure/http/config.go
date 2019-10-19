package http

import "time"

//Config contains a HTTP adapter configuration.
type Config struct {
	Address     string        `validate:"required"`
	ReadTimeout time.Duration `validate:"required"`
}
