package application

import "errors"

var (
	ErrNotFound             = errors.New("fuel record not found")
	ErrInvalidFuelData      = errors.New("invalid fuel data")
	ErrInvalidFuelKind      = errors.New("invalid fuel type")
	ErrInvalidStatus        = errors.New("invalid fuel transaction status")
	ErrCompletedImmutable   = errors.New("completed fuel transactions cannot be modified directly")
	ErrInvalidTransition    = errors.New("invalid fuel transaction transition")
	ErrDuplicateCNPJ        = errors.New("fuel station cnpj already exists for tenant")
	ErrAdjustmentRequired   = errors.New("fuel transaction corrections require an adjustment")
	ErrTenantIDRequired     = errors.New("tenant id is required")
	ErrReasonRequired       = errors.New("reason is required")
	ErrTransactionRequired  = errors.New("transaction id is required")
	ErrFuelTransactionFinal = errors.New("final fuel transaction cannot be deleted")
)
