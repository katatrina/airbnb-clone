// Package apperror provides application-level error types that bridge
// the gap between internal errors (database, external services) and
// HTTP responses. This allows the service layer to communicate errors
// to the handler layer in a structured way, without leaking implementation details.
package apperror

import (
	"fmt"
	"net/http"
)

// ErrorType categorizes errors to help handlers determine the appropriate HTTP status code.
// This is the key insight: internal errors (like pgx.ErrNoRows) get translated into
// semantic error types (like ErrorTypeNotFound) that handlers can understand without
// knowing about the underlying database.
type ErrorType int

const (
	// Client errors - caused by invalid input or unauthorized access
	ErrorTypeValidation   ErrorType = iota // 400 Bad Request
	ErrorTypeUnauthorized                  // 401 Unauthorized
	ErrorTypeForbidden                     // 403 Forbidden
	ErrorTypeNotFound                      // 404 Not Found
	ErrorTypeConflict                      // 409 Conflict

	// Server errors - caused by internal failures
	ErrorTypeInternal // 500 Internal Server Error
)

// AppError is the standard error type used throughout the application.
// It carries enough information for:
// 1. Handlers to generate appropriate HTTP responses (Type, Code, Message)
// 2. Logging systems to capture debug information (Detail, Err)
// 3. Clients to handle errors programmatically (Code, Fields)
type AppError struct {
	Type    ErrorType    // Determines HTTP status code
	Code    int          // Business error code for client-side handling
	Message string       // User-facing message (safe to return to client)
	Detail  string       // Internal detail for logging (NOT returned to client)
	Err     error        // Original error for debugging (NOT returned to client)
	Fields  []FieldError // Validation errors for specific fields
}

// FieldError represents a validation error for a specific field.
// This allows the client to highlight exactly which form fields have problems.
type FieldError struct {
	Field   string `json:"field"`   // Field name (e.g., "email", "password")
	Code    string `json:"code"`    // Machine-readable code (e.g., "REQUIRED", "INVALID_FORMAT")
	Message string `json:"message"` // Human-readable message
}

// Error implements the error interface.
// The message includes the wrapped error if present, useful for logging.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error, enabling errors.Is() and errors.As() to work.
// This is important for error chain inspection in Go 1.13+.
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus maps the error type to an HTTP status code.
// This centralizes the mapping logic so handlers don't need to know the details.
func (e *AppError) HTTPStatus() int {
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
