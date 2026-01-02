package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============ Success responses ============

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, New().Success(data).Build())
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, New().Success(data).Build())
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ============ Error responses ============

func BadRequest(c *gin.Context, code int, message string) {
	c.JSON(http.StatusBadRequest, New().Error(code, message).Build())
}

func BadRequestWithErrors(c *gin.Context, code int, message string, errors []FieldError) {
	c.JSON(http.StatusBadRequest, New().Error(code, message).WithErrors(errors).Build())
}

func Unauthorized(c *gin.Context, message string) {
	// 4010 = Authentication required
	// 4011 = Token invalid or expired
	c.JSON(http.StatusUnauthorized, New().Error(4010, message).Build())
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, New().Error(4030, message).Build())
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, New().Error(4040, message).Build())
}

func Conflict(c *gin.Context, code int, message string) {
	c.JSON(http.StatusConflict, New().Error(code, message).Build())
	http.NewFileTransport()
}

func InternalError(c *gin.Context) {
	// Never send internal server error detail to client
	c.JSON(http.StatusInternalServerError, New().Error(5000, "Internal server error").Build())
}

// ============ Paginated response ============

func OKWithPagination(c *gin.Context, data any, page, perPage int, total int64) {
	c.JSON(http.StatusOK, New().Success(data).WithPagination(page, perPage, total).Build())
}
