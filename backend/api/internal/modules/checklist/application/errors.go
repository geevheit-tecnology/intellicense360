package application

import "errors"

var (
	ErrNotFound                  = errors.New("checklist resource not found")
	ErrValidation                = errors.New("checklist validation failed")
	ErrInvalidStatus             = errors.New("invalid checklist status transition")
	ErrInvalidAnswer             = errors.New("invalid checklist answer")
	ErrPublishedVersionImmutable = errors.New("published checklist template version cannot be modified")
	ErrRequiredItemsUnanswered   = errors.New("required checklist items are unanswered")
	ErrRequiredEvidenceMissing   = errors.New("required checklist evidence is missing")
	ErrRequiredSignatureMissing  = errors.New("required checklist signature is missing")
	ErrInvalidTransition         = errors.New("invalid checklist transition")
)
