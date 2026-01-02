package response

const TestVal = 111

// Response is the standard structure for all API responses.
type Response struct {
	Success bool         `json:"success"`
	Code    int          `json:"code"` // Business code, 0 = success
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"` // Validation errors
	Meta    *Meta        `json:"meta,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	RequestID  string      `json:"requestId"`
	Timestamp  int64       `json:"timestamp"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
