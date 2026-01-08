// Package response provides a standardized way to format API responses.
// All endpoints should use this package to ensure consistent response structure.
//
// Response structure:
//
//	{
//	  "success": true/false,
//	  "code": "OK",           // Business error code
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
	"github.com/katatrina/airbnb-clone/pkg/validator"
)

// Response is the standard structure for all API responses.
type Response struct {
	Success bool                   `json:"success"`
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Data    any                    `json:"data,omitempty"`
	Meta    *Meta                  `json:"meta,omitempty"`
	Errors  []validator.FieldError `json:"errors,omitempty"`
}

// Meta contains metadata about the response.
type Meta struct {
	RequestID  string      `json:"requestId"`
	Timestamp  int64       `json:"timestamp"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination metadata.
type Pagination struct {
	Page       int   `json:"page"`       // Current page number (1-indexed)
	PerPage    int   `json:"perPage"`    // Items per page
	Total      int64 `json:"total"`      // Total number of items
	TotalPages int   `json:"totalPages"` // Total number of pages
}

// Builder constructs a Response step by step.
// This makes the code more readable and allows method chaining.
//
// Example:
//
//	New().Success(data).WithPagination(1, 10, 100).Build()
type Builder struct {
	resp Response
}

// New creates a new response builder with default metadata.
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
func (b *Builder) WithRequestID(id string) *Builder {
	b.resp.Meta.RequestID = id
	return b
}

// Success marks the response as successful and sets the data.
func (b *Builder) Success(data any) *Builder {
	b.resp.Success = true
	b.resp.Code = CodeSuccess
	b.resp.Message = "OK"
	b.resp.Data = data
	return b
}

// Error marks the response as failed with an error code and message.
func (b *Builder) Error(code ErrorCode, message string) *Builder {
	b.resp.Success = false
	b.resp.Code = code
	b.resp.Message = message
	return b
}

// WithErrors adds field-level validation errors.
func (b *Builder) WithErrors(errors []validator.FieldError) *Builder {
	b.resp.Errors = errors
	return b
}

// WithPagination adds pagination metadata.
func (b *Builder) WithPagination(page, perPage int, total int64) *Builder {
	// TODO: understand pagination logic
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
