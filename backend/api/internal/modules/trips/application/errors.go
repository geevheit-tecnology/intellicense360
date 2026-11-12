package application

import "errors"

var (
	ErrNotFound            = errors.New("trip not found")
	ErrValidation          = errors.New("invalid trip data")
	ErrInvalidStatusChange = errors.New("invalid trip status transition")
	ErrOdometerRegression  = errors.New("end odometer reading must be greater than or equal to start odometer")
	ErrVehicleRequired     = errors.New("vehicle is required for dispatch")
	ErrDriverRequired      = errors.New("driver is required for dispatch")
)

