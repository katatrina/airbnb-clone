package validator

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	dangerousChars = regexp.MustCompile(`[<>"';&]`)
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
			message: "{0} must contain uppercase, lowercase, and number",
		},
		{
			tag:     "safename",
			fn:      validateSafeName,
			message: "{0} contains invalid characters",
		},
	}

	for _, r := range rules {
		_ = validate.RegisterValidation(r.tag, r.fn)
		registerTranslation(r.tag, r.message)
	}
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var hasUpper, hasLower, hasNumber bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}

func validateSafeName(fl validator.FieldLevel) bool {
	return !dangerousChars.MatchString(fl.Field().String())
}
