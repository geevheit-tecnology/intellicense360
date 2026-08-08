package application

import "errors"

var (
	ErrNotFound      = errors.New("maintenance resource not found")
	ErrValidation    = errors.New("maintenance validation failed")
	ErrInvalidStatus = errors.New("invalid maintenance status transition")
)
