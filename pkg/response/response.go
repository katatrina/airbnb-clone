// Package response provides a standardized way to format API responses.
// All endpoints should use this package to ensure consistent response structure.
//
// Response structure:
//
//	{
//	  "success": true/false,
//	  "code": 0,              // 0 for success, error code for failures
//	  "message": "OK",        // Human-readable message
//	  "data": {...},          // Payload (only for success)
//	  "errors": [...],        // Field errors (only for validation failures)
//	  "meta": {
//	    "requestId": "...",   // For tracing/debugging
//	    "timestamp": 123456,  // Server time
//	    "pagination": {...}   // For list endpoints
//	  }
//	}
package response

import (
	"time"

	"github.com/google/uuid"
)

// Response is the standard structure for all API responses.
// Using a consistent structure makes it easier for frontend to handle responses.
type Response struct {
	// Success indicates whether the request was successful.
	// Frontend can check this first before looking at other fields.
	Success bool `json:"success"`

	// Code is a machine-readable error code.
	// 0 = success, non-zero = specific error type.
	// This allows frontend to handle errors programmatically.
	Code int `json:"code"`

	// Message is a human-readable description.
	// Can be displayed to users or used for debugging.
	Message string `json:"message"`

	// Data contains the response payload.
	// Only present for successful responses.
	Data any `json:"data,omitempty"`

	// Meta contains metadata about the request/response.
	// Useful for debugging, tracing, and pagination.
	Meta *Meta `json:"meta,omitempty"`

	// Errors contains field-level validation errors.
	// Only present when there are validation failures.
	Errors []FieldError `json:"errors,omitempty"`
}

// FieldError represents a validation error for a specific field.
// This allows frontend to highlight the exact form field that has an error.
type FieldError struct {
	Field   string `json:"field"`   // e.g., "email", "password"
	Code    string `json:"code"`    // e.g., "REQUIRED", "INVALID_FORMAT"
	Message string `json:"message"` // e.g., "Email is required"
}

// Meta contains metadata about the response.
type Meta struct {
	// RequestID is a unique identifier for this request.
	// Use this to trace logs and debug issues in production.
	RequestID string `json:"requestId"`

	// Timestamp is when the response was generated (Unix seconds).
	// Using Unix timestamp makes it easy for any frontend to parse.
	Timestamp int64 `json:"timestamp"`

	// Pagination contains pagination info for list endpoints.
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination metadata.
type Pagination struct {
	Page       int   `json:"page"`       // Current page number (1-indexed)
	PerPage    int   `json:"perPage"`    // Items per page
	Total      int64 `json:"total"`      // Total number of items
	TotalPages int   `json:"totalPages"` // Total number of pages
}

// ============================================================================
// Builder Pattern
// ============================================================================
// We use the Builder pattern to construct responses.
// This makes the code more readable and allows method chaining.
//
// Example:
//   New().Success(data).WithPagination(1, 10, 100).Build()

// Builder constructs a Response step by step.
type Builder struct {
	resp Response
}

// New creates a new response builder with default metadata.
// Every response gets a unique request ID and timestamp automatically.
func New() *Builder {
	return &Builder{
		resp: Response{
			Meta: &Meta{
				RequestID: uuid.NewString(),
				Timestamp: time.Now().Unix(),
			},
		},
	}
}

// WithRequestID sets a custom request ID.
// Use this when you have a request ID from middleware/tracing.
func (b *Builder) WithRequestID(id string) *Builder {
	b.resp.Meta.RequestID = id
	return b
}

// Success marks the response as successful and sets the data.
func (b *Builder) Success(data any) *Builder {
	b.resp.Success = true
	b.resp.Code = 0
	b.resp.Message = "OK"
	b.resp.Data = data
	return b
}

// Error marks the response as failed with an error code and message.
func (b *Builder) Error(code int, message string) *Builder {
	b.resp.Success = false
	b.resp.Code = code
	b.resp.Message = message
	return b
}

// WithErrors adds field-level validation errors.
func (b *Builder) WithErrors(errors []FieldError) *Builder {
	b.resp.Errors = errors
	return b
}

// WithPagination adds pagination metadata.
// TotalPages is calculated automatically from total and perPage.
func (b *Builder) WithPagination(page, perPage int, total int64) *Builder {
	// Calculate total pages using ceiling division
	// Example: 25 items with 10 per page = 3 pages
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	b.resp.Meta.Pagination = &Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
	return b
}

// Build returns the constructed Response.
func (b *Builder) Build() Response {
	return b.resp
}
