package application

import "errors"

var (
	ErrNotFound         = errors.New("intelligence record not found")
	ErrInvalidData      = errors.New("invalid intelligence data")
	ErrInsufficientData = errors.New("insufficient data")
	ErrDuplicateInsight = errors.New("duplicate insight")
)
