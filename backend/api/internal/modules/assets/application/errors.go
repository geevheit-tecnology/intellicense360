package application

import "errors"

var (
	ErrNotFound          = errors.New("asset resource not found")
	ErrValidation        = errors.New("asset validation failed")
	ErrInternalCodeTaken = errors.New("internal code already exists")
	ErrSerialNumberTaken = errors.New("serial number already exists")
	ErrAssetTagTaken     = errors.New("asset tag already exists")
	ErrInvalidStatus     = errors.New("invalid asset status")
)
