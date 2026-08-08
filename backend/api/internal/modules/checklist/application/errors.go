package application

import "errors"

var (
	ErrNotFound      = errors.New("checklist resource not found")
	ErrValidation    = errors.New("checklist validation failed")
	ErrInvalidStatus = errors.New("invalid checklist status transition")
	ErrInvalidAnswer = errors.New("invalid checklist answer")
)
