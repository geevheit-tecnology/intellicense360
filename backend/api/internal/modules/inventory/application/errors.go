package application

import "errors"

var (
	ErrNotFound          = errors.New("inventory resource not found")
	ErrValidation        = errors.New("invalid inventory data")
	ErrSKUTaken          = errors.New("sku already exists")
	ErrInternalCodeTaken = errors.New("internal code already exists")
	ErrInvalidUnit       = errors.New("invalid unit of measure")
	ErrInvalidStockRange = errors.New("minimum stock cannot be greater than maximum stock")
)
