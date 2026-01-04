package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============ Success responses ============

// OK sends a 200 response with data.
// Use for successful GET requests or operations that return data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, New().Success(data).Build())
}

// Created sends a 201 response with data.
// Use after successfully creating a new resource (POST).
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, New().Success(data).Build())
}

// NoContent sends a 204 response with no body.
// Use after successful DELETE or updates that don't return data.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ============ Error responses ============

// BadRequest sends a 400 response.
// Use when the client sends invalid data (validation errors, malformed JSON).
func BadRequest(c *gin.Context, code int, message string) {
	c.JSON(http.StatusBadRequest, New().Error(code, message).Build())
}

// BadRequestWithErrors sends a 400 response with field-level errors.
// Use for validation errors that need to highlight specific fields.
func BadRequestWithErrors(c *gin.Context, code int, message string, errors []FieldError) {
	c.JSON(http.StatusBadRequest, New().Error(code, message).WithErrors(errors).Build())
}

// Unauthorized sends a 401 response.
// Use when authentication is required but missing or invalid.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, New().Error(CodeUnauthorized, message).Build())
}

// Forbidden sends a 403 response.
// Use when the user is authenticated but doesn't have permission.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, New().Error(CodeForbidden, message).Build())
}

// NotFound sends a 404 response.
// Use when the requested resource doesn't exist.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, New().Error(CodeNotFound, message).Build())
}

// Conflict sends a 409 response.
// Use when there's a conflict with existing data (e.g., duplicate email).
func Conflict(c *gin.Context, code int, message string) {
	c.JSON(http.StatusConflict, New().Error(code, message).Build())
}

// InternalError sends a 500 response.
// Use for unexpected server errors. Never expose internal details to client.
func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, New().Error(CodeInternalError, "Internal server error").Build())
}

// ============ Paginated response ============

// OKWithPagination sends a 200 response with pagination metadata.
// Use for list endpoints that support pagination.
func OKWithPagination(c *gin.Context, data any, page, perPage int, total int64) {
	c.JSON(http.StatusOK, New().Success(data).WithPagination(page, perPage, total).Build())
}
