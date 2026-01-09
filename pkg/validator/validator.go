// Package validator provides request validation and normalization.
package validator

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

// FieldErrorCode represents the type of validation error.
type FieldErrorCode string

const (
	FieldCodeRequired      FieldErrorCode = "REQUIRED"
	FieldCodeInvalidFormat FieldErrorCode = "INVALID_FORMAT"
	FieldCodeTooShort      FieldErrorCode = "TOO_SHORT"
	FieldCodeTooLong       FieldErrorCode = "TOO_LONG"
	FieldCodeMinValue      FieldErrorCode = "MIN_VALUE"
	FieldCodeMaxValue      FieldErrorCode = "MAX_VALUE"
)

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field   string         `json:"field"`
	Value   any            `json:"value"`
	Code    FieldErrorCode `json:"code"`
	Message string         `json:"message"`
}

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("failed to get validator engine")
	}
	validate = v

	// Setup translator
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(v, trans)

	// Use JSON tag name
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return field.Name
		}
		return name
	})

	registerMessages()
	registerRules()
}

// TranslateValidationErrors converts validator.ValidationErrors to a slice of FieldError.
func TranslateValidationErrors(err validator.ValidationErrors) []FieldError {
	var fieldErrors []FieldError
	for _, e := range err {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   toCamelCase(e.Field()),
			Value:   e.Value(),
			Code:    mapValidationTag(e.Tag()),
			Message: e.Translate(trans),
		})
	}

	return fieldErrors
}

// mapValidationTag maps validator tags to FieldErrorCode constants.
func mapValidationTag(tag string) FieldErrorCode {
	switch tag {
	case "required":
		return FieldCodeRequired
	case "email":
		return FieldCodeInvalidFormat
	case "min":
		return FieldCodeTooShort
	case "max", "maxbytes":
		return FieldCodeTooLong
	case "gte":
		return FieldCodeMinValue
	case "lte":
		return FieldCodeMaxValue
	default:
		return FieldErrorCode(strings.ToUpper(tag))
	}
}

// toCamelCase converts PascalCase to camelCase.
// DisplayName -> displayName
// UserID -> userId (properly handles acronyms)
func toCamelCase(s string) string {
	if s == "" {
		return ""
	}

	// Convert first character to lowercase
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])

	// Handle acronyms: UserID -> userId, not userID
	// Find where the acronym ends (when we hit a lowercase letter)
	for i := 1; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i]) && unicode.IsLower(runes[i+1]) {
			// This is where acronym ends
			// e.g., "UserID" -> i=4 (D), so convert I to lowercase
			// Result: "userId"
			break
		}
		if unicode.IsUpper(runes[i]) {
			runes[i] = unicode.ToLower(runes[i])
		} else {
			break
		}
	}

	return string(runes)
}
