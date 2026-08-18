package application

import "errors"

var (
	ErrCommandItemNotFound     = errors.New("command item not found")
	ErrInvalidCommandItemType  = errors.New("invalid command item type")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidSeverity         = errors.New("invalid severity")
	ErrInvalidPriority         = errors.New("invalid priority")
	ErrInvalidRiskScore        = errors.New("invalid risk score")
	ErrInvalidImpactScore      = errors.New("invalid impact score")
	ErrInvalidConfidence       = errors.New("invalid confidence")
	ErrDuplicateCommandItem    = errors.New("duplicate command item")
	ErrTenantMismatch          = errors.New("tenant mismatch")
	ErrCommandActionNotAllowed = errors.New("command action not allowed")
	ErrInvalidSLA              = errors.New("invalid sla")
	ErrInvalidData             = errors.New("invalid mission control data")
)
