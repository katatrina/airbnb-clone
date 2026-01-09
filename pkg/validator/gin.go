package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ShouldBindJSON binds JSON request body to obj, normalizes it, then validates.
// This is a Gin-specific adapter over the framework-agnostic validation logic.
//
// This is a drop-in replacement for gin.Context.ShouldBindJSON() with auto-normalization.
//
// Usage:
//
//	var req RegisterRequest
//	if err := validator.ShouldBindJSON(c, &req); err != nil {
//	    response.HandleJSONBindingError(c, err)
//	    return
//	}
//
// Note: This function is Gin-specific. For other frameworks, implement a similar
// adapter that calls NormalizeStruct() and validates using the standard validator.
func ShouldBindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindBodyWith(obj, binding.JSON); err != nil {
		return err
	}

	NormalizeStruct(obj)

	if err := binding.Validator.ValidateStruct(obj); err != nil {
		return err
	}

	return nil
}
