package response

import (
	"time"

	"github.com/google/uuid"
)

// Response is the standard structure for all API responses.
type Response struct {
	Success bool `json:"success"` // Quickly know whether the request is successful or not, no need to check HTTP status code

	// Business error code, allowing coarse-grained (request-level) error handling
	//
	// 0 = success, non-zero = error
	Code    int          `json:"code"`
	Message string       `json:"message"`          // Human-readable message, used to debug or sometimes display for user
	Data    any          `json:"data,omitempty"`   // Payload
	Meta    *Meta        `json:"meta,omitempty"`   // Metadata about request/response, not business data
	Errors  []FieldError `json:"errors,omitempty"` // Validation errors
}

type FieldError struct {
	Field string `json:"field"` // e.g. email, password, phone,...

	// Machine-readable error code, SNAKE_CASE, for client mapping
	//
	// Allow fine-grained (field-level) error handling
	Code    string `json:"code"`
	Message string `json:"message"` // Human-readable validation error message
}

type Meta struct {
	RequestID  string      `json:"requestId"` // Quickly identify the request, help debugging (trace log) extremely fast, especially in production
	Timestamp  int64       `json:"timestamp"` // Server time when response is created, with Unix universal timestamp, easy for frontend to parse
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`       // Current page
	PerPage    int   `json:"perPage"`    // Number of items per page
	Total      int64 `json:"total"`      // Total items (use int64 as it can be very large)
	TotalPages int   `json:"totalPages"` // Total pages
}

type Builder struct {
	resp Response
}

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

// WithRequestID is used to set request ID from middleware (for tracing)
func (b *Builder) WithRequestID(id string) *Builder {
	b.resp.Meta.RequestID = id
	return b
}

func (b *Builder) Success(data any) *Builder {
	b.resp.Success = true
	b.resp.Code = 0
	b.resp.Message = "OK"
	b.resp.Data = data
	return b
}

func (b *Builder) Error(code int, message string) *Builder {
	b.resp.Success = false
	b.resp.Code = code
	b.resp.Message = message
	return b
}

func (b *Builder) WithErrors(errors []FieldError) *Builder {
	b.resp.Errors = errors
	return b
}

func (b *Builder) WithPagination(page, perPage int, total int64) *Builder {
	// TODO: Explain pagination calculation logic here
	b.resp.Meta.Pagination = &Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: (int(total) + perPage - 1) / perPage,
	}
	return b
}

func (b *Builder) Build() Response {
	return b.resp
}
