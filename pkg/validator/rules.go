package validator

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func registerRules() {
	rules := []struct {
		tag     string
		fn      validator.Func
		message string
	}{
		{
			tag:     "strongpass",
			fn:      validateStrongPassword,
			message: "{0} must contain at least one uppercase, one lowercase, one number, and one special character",
		},
		{
			tag:     "displayname",
			fn:      validateDisplayName,
			message: "{0} must contain at least one letter and must not contain URLs",
		},
		{
			tag:     "maxbytes",
			fn:      ValidateMaxBytes,
			message: "{0} is too long",
		},
	}

	for _, r := range rules {
		_ = validate.RegisterValidation(r.tag, r.fn)
		registerTranslation(r.tag, r.message)
	}
}

// validateStrongPassword: At least 8 chars, has uppercase, lowercase, number
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateDisplayName: Must contain at least one letter, no URLs
func validateDisplayName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	// Must have at least one letter
	if !regexp.MustCompile(`\p{L}`).MatchString(name) {
		return false
	}

	// No URLs (spam prevention)
	if regexp.MustCompile(`https?://`).MatchString(strings.ToLower(name)) {
		return false
	}

	return true
}

// ValidateMaxBytes checks that the string length in bytes does not exceed the specified limit.
func ValidateMaxBytes(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	param := fl.Param()

	limit, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	return len(field) <= limit
}
