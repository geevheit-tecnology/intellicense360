package events

import "errors"

var (
	ErrInvalidEvent       = errors.New("invalid domain event")
	ErrInvalidEventType   = errors.New("invalid event type")
	ErrHandlerFailed      = errors.New("event handler failed")
	ErrDuplicateEvent     = errors.New("duplicate event")
	ErrDuplicateConsumer  = errors.New("event already processed by consumer")
	ErrOutboxNotFound     = errors.New("outbox record not found")
	ErrRetryLimitExceeded = errors.New("event retry limit exceeded")
)
