package application

import "errors"

var (
	ErrConflict                   = errors.New("fleet resource conflict")
	ErrNotFound                   = errors.New("fleet resource not found")
	ErrValidation                 = errors.New("fleet validation failed")
	ErrLicensePlateTaken          = errors.New("license plate already exists")
	ErrChassisTaken               = errors.New("chassis already exists")
	ErrRenavamTaken               = errors.New("renavam already exists")
	ErrOdometerRegression         = errors.New("odometer reading cannot be lower than previous reading without valid correction")
	ErrOdometerCorrectionRequired = errors.New("odometer correction reason must have at least 10 characters")
)
