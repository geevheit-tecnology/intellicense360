package application

import "errors"

var (
	ErrNotFound              = errors.New("financial record not found")
	ErrInvalidFinancialData  = errors.New("invalid financial data")
	ErrInvalidStatus         = errors.New("invalid financial status")
	ErrInvalidTransition     = errors.New("invalid financial transition")
	ErrFinalizedImmutable    = errors.New("paid or received financial transactions cannot be modified directly")
	ErrClosedPeriodImmutable = errors.New("closed financial period does not allow transaction mutation")
	ErrAdjustmentRequired    = errors.New("financial corrections require an adjustment")
)
