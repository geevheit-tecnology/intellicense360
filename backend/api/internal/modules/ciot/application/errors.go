package application

import "errors"

var (
	ErrNotFound                  = errors.New("ciot not found")
	ErrInvalidData               = errors.New("invalid ciot data")
	ErrInvalidTransition         = errors.New("invalid ciot status transition")
	ErrDuplicateRequest          = errors.New("duplicate ciot idempotency request")
	ErrFinalizedImmutable        = errors.New("finalized ciot is immutable")
	ErrExternalReferenceConflict = errors.New("external ciot reference conflict")
)
