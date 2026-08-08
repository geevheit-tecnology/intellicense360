package application

import "errors"

var (
	ErrNotFound          = errors.New("tire resource not found")
	ErrValidation        = errors.New("tire validation failed")
	ErrSerialNumberTaken = errors.New("serial number already exists")
	ErrFireNumberTaken   = errors.New("fire number already exists")
	ErrDisposedTire      = errors.New("disposed tire cannot be moved")
	ErrInvalidStatus     = errors.New("invalid tire status")
	ErrInvalidCondition  = errors.New("invalid tire condition")
	ErrInvalidTransition = errors.New("invalid tire lifecycle transition")
)
