package application

import "errors"

var (
	ErrNotFound          = errors.New("fleet resource not found")
	ErrValidation        = errors.New("fleet validation failed")
	ErrLicensePlateTaken = errors.New("license plate already exists")
	ErrChassisTaken      = errors.New("chassis already exists")
	ErrRenavamTaken      = errors.New("renavam already exists")
)
