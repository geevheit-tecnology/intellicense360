package application

import "errors"

var (
	ErrNotFound          = errors.New("supplier resource not found")
	ErrValidation        = errors.New("invalid supplier data")
	ErrCNPJTaken         = errors.New("cnpj already exists")
	ErrInvalidCNPJ       = errors.New("invalid cnpj")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidState      = errors.New("invalid state")
	ErrInvalidStatus     = errors.New("invalid supplier status")
	ErrInvalidType       = errors.New("invalid supplier type")
	ErrInvalidTransition = errors.New("invalid supplier status transition")
)
