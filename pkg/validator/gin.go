package validator

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

// ShouldBindJSON binds JSON request body to obj, normalizes it, then validates.
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
func ShouldBindJSON(c *gin.Context, obj interface{}) error {
	// Read body and allow re-reading
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Bind JSON without validation
	if err := json.Unmarshal(body, obj); err != nil {
		return err
	}

	// Normalize first
	NormalizeStruct(obj)

	// Then validate (now with normalized values)
	if err = validate.Struct(obj); err != nil {
		return err
	}

	return nil
}
