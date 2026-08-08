package errors

import "net/http"

type Code string

const (
	CodeInternal      Code = "internal_error"
	CodeUnauthorized  Code = "unauthorized"
	CodeForbidden     Code = "forbidden"
	CodeInvalidInput  Code = "invalid_input"
	CodeNotFound      Code = "not_found"
	CodeTenantMissing Code = "tenant_missing"
)

type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e AppError) Error() string {
	return e.Message
}

func New(code Code, message string, status int) AppError {
	return AppError{Code: code, Message: message, Status: status}
}

func Internal(message string) AppError {
	return New(CodeInternal, message, http.StatusInternalServerError)
}

func Unauthorized(message string) AppError {
	return New(CodeUnauthorized, message, http.StatusUnauthorized)
}

func TenantMissing() AppError {
	return New(CodeTenantMissing, "tenant header is required", http.StatusBadRequest)
}
