package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is our standard error type that carries rich context about what went wrong.
// This replaces bare error returns throughout the app with structured errors
// that include HTTP status codes, error codes for clients, and detailed messages.
//
// Why we need this:
// - Domain errors (like "invalid credentials") need to map to HTTP status codes
// - Clients need machine-readable error codes to handle errors programmatically
// - Developers need detailed messages for debugging
// - We want to hide internal errors from clients while logging them
type AppError struct {
	// Code is a machine-readable error code for clients
	// Examples: "AUTH_INVALID_CREDENTIALS", "USER_NOT_FOUND"
	Code string `json:"code"`

	// Message is a human-readable error message safe to show to clients
	Message string `json:"message"`

	// Details provides additional context about the error
	// This might include validation errors, field names, etc.
	Details map[string]interface{} `json:"details,omitempty"`

	// StatusCode is the HTTP status code to return
	StatusCode int `json:"-"`

	// Err is the underlying error that caused this
	// We don't expose this to clients but use it for logging
	Err error `json:"-"`
}

// Error implements the error interface so AppError can be used anywhere error is expected
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work with AppError
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatusCode exposes the HTTP status code for this error.
// This makes AppError compatible with httputil responders.
func (e *AppError) HTTPStatusCode() int {
	return e.StatusCode
}

// New creates a new AppError with all fields specified
func New(code string, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
		Details:    make(map[string]interface{}),
	}
}

// WithDetails adds additional context to an AppError
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithError wraps an underlying error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts the AppError from an error chain
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// Common error codes - define these as constants for consistency
const (
	ErrCodeInvalidCredentials     = "AUTH_INVALID_CREDENTIALS"
	ErrCodeEmailNotVerified       = "AUTH_EMAIL_NOT_VERIFIED"
	ErrCodeSessionExpired         = "AUTH_SESSION_EXPIRED"
	ErrCodeSessionBlocked         = "AUTH_SESSION_BLOCKED"
	ErrCodeInvalidToken           = "AUTH_INVALID_TOKEN"
	ErrCodeTokenExpired           = "AUTH_TOKEN_EXPIRED"
	ErrCodeUnauthorized           = "AUTH_UNAUTHORIZED"
	ErrCodeInsufficientPermission = "AUTH_INSUFFICIENT_PERMISSION"

	ErrCodeNotFound      = "RESOURCE_NOT_FOUND"
	ErrCodeAlreadyExists = "RESOURCE_ALREADY_EXISTS"
	ErrCodeConflict      = "RESOURCE_CONFLICT"

	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeInvalidInput     = "VALIDATION_INVALID_INPUT"
	ErrCodeMissingField     = "VALIDATION_MISSING_FIELD"

	ErrCodeInternalError        = "INTERNAL_ERROR"
	ErrCodeDatabaseError        = "INTERNAL_DATABASE_ERROR"
	ErrCodeExternalServiceError = "INTERNAL_EXTERNAL_SERVICE_ERROR"

	ErrCodeBusinessRuleViolation = "BUSINESS_RULE_VIOLATION"
	ErrCodeOperationNotAllowed   = "BUSINESS_OPERATION_NOT_ALLOWED"
)

// Pre-defined common errors for consistency
// These can be reused throughout the app

func InvalidCredentials(err error) *AppError {
	return New(
		ErrCodeInvalidCredentials,
		"Invalid email or password",
		http.StatusUnauthorized,
		err,
	)
}

func EmailNotVerified() *AppError {
	return New(
		ErrCodeEmailNotVerified,
		"Please verify your email address",
		http.StatusForbidden,
		nil,
	)
}

func SessionExpired(err error) *AppError {
	return New(
		ErrCodeSessionExpired,
		"Your session has expired. Please log in again",
		http.StatusUnauthorized,
		err,
	)
}

func SessionBlocked(reason string) *AppError {
	return New(
		ErrCodeSessionBlocked,
		"Your session has been blocked",
		http.StatusForbidden,
		nil,
	).WithDetails("reason", reason)
}

func InvalidToken(err error) *AppError {
	return New(
		ErrCodeInvalidToken,
		"Invalid or malformed token",
		http.StatusUnauthorized,
		err,
	)
}

func TokenExpired(err error) *AppError {
	return New(
		ErrCodeTokenExpired,
		"Token has expired",
		http.StatusUnauthorized,
		err,
	)
}

func Unauthorized(message string) *AppError {
	return New(
		ErrCodeUnauthorized,
		message,
		http.StatusUnauthorized,
		nil,
	)
}

func NotFound(resource string, identifier string) *AppError {
	return New(
		ErrCodeNotFound,
		fmt.Sprintf("%s not found", resource),
		http.StatusNotFound,
		nil,
	).WithDetails("resource", resource).WithDetails("identifier", identifier)
}

func AlreadyExists(resource string, identifier string) *AppError {
	return New(
		ErrCodeAlreadyExists,
		fmt.Sprintf("%s already exists", resource),
		http.StatusConflict,
		nil,
	).WithDetails("resource", resource).WithDetails("identifier", identifier)
}

func ValidationFailed(message string) *AppError {
	return New(
		ErrCodeValidationFailed,
		message,
		http.StatusBadRequest,
		nil,
	)
}

func ValidationFailedWithDetails(message string, details map[string]interface{}) *AppError {
	err := New(
		ErrCodeValidationFailed,
		message,
		http.StatusBadRequest,
		nil,
	)
	err.Details = details
	return err
}

func InvalidInput(field string, reason string) *AppError {
	return New(
		ErrCodeInvalidInput,
		fmt.Sprintf("Invalid input for field: %s", field),
		http.StatusBadRequest,
		nil,
	).WithDetails("field", field).WithDetails("reason", reason)
}

func InternalError(err error) *AppError {
	return New(
		ErrCodeInternalError,
		"An internal error occurred. Please try again later",
		http.StatusInternalServerError,
		err,
	)
}

func DatabaseError(operation string, err error) *AppError {
	return New(
		ErrCodeDatabaseError,
		"A database error occurred",
		http.StatusInternalServerError,
		err,
	).WithDetails("operation", operation)
}

func BusinessRuleViolation(rule string, message string) *AppError {
	return New(
		ErrCodeBusinessRuleViolation,
		message,
		http.StatusUnprocessableEntity,
		nil,
	).WithDetails("rule", rule)
}

func OperationNotAllowed(operation string, reason string) *AppError {
	return New(
		ErrCodeOperationNotAllowed,
		fmt.Sprintf("Operation not allowed: %s", operation),
		http.StatusForbidden,
		nil,
	).WithDetails("operation", operation).WithDetails("reason", reason)
}
