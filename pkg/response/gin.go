package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	validatorV10 "github.com/go-playground/validator/v10"
	"github.com/katatrina/airbnb-clone/pkg/validator"
)

// ============ Success responses ============

// OK sends a 200 response with data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, New().Success(data).Build())
}

// Created sends a 201 response with data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, New().Success(data).Build())
}

// OKWithPagination sends a 200 response with pagination metadata.
func OKWithPagination(c *gin.Context, data any, page, perPage int, total int64) {
	c.JSON(http.StatusOK, New().Success(data).WithPagination(page, perPage, total).Build())
}

// NoContent sends a 204 response with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ============ Error responses ============

// BadRequest sends a 400 response.
func BadRequest(c *gin.Context, code ErrorCode, message string) {
	c.JSON(http.StatusBadRequest, New().Error(code, message).Build())
}

// BadRequestWithErrors sends a 400 response with field-level errors.
func BadRequestWithErrors(c *gin.Context, code ErrorCode, message string, errors []validator.FieldError) {
	c.JSON(http.StatusBadRequest, New().Error(code, message).WithErrors(errors).Build())
}

// Unauthorized sends a 401 response.
func Unauthorized(c *gin.Context, code ErrorCode, message string) {
	c.JSON(http.StatusUnauthorized, New().Error(code, message).Build())
}

// NotFound sends a 404 response.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, New().Error(CodeResourceNotFound, message).Build())
}

// Conflict sends a 409 response.
func Conflict(c *gin.Context, code ErrorCode, message string) {
	c.JSON(http.StatusConflict, New().Error(code, message).Build())
}

// InternalServerError sends a 500 response.
// Never expose internal details to client.
func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError,
		New().Error(CodeInternalServerError, "Internal server error. Please try again later.").Build(),
	)
}

// ============ Helper ============

// HandleJSONBindingError properly handles different types of validator.ShouldBindJSON errors.
// It distinguishes between JSON parsing errors and validation errors,
// returning appropriate error codes and messages.
//
// Usage:
//
//	if err := validator.ShouldBindJSON(&req); err != nil {
//	    response.HandleJSONBindingError(c, err)
//	    return
//	}
func HandleJSONBindingError(c *gin.Context, err error) {
	// Validation errors
	var validationErrors validatorV10.ValidationErrors
	if errors.As(err, &validationErrors) {
		fieldErrors := validator.TranslateValidationErrors(validationErrors)
		BadRequestWithErrors(c, CodeValidationFailed, "Validation failed", fieldErrors)
		return
	}

	// JSON parsing errors or other errors
	BadRequest(c, CodeJSONFormatInvalid, "Request body must be valid JSON")
}
