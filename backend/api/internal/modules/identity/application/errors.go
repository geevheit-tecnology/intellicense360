package application

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrConflict           = errors.New("identity resource conflict")
	ErrValidation         = errors.New("identity validation failed")
	ErrNotFound           = errors.New("identity resource not found")
	ErrForbidden          = errors.New("forbidden")
	ErrPasswordPolicy     = errors.New("password does not satisfy policy")
)
