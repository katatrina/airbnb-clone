package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func registerMessages() {
	messages := map[string]string{
		"required": "{0} is required",
		"email":    "{0} must be a valid email address",
		"min":      "{0} must be at least {1} characters",
		"max":      "{0} must be at most {1} characters",
		"gte":      "{0} must be at least {1}",
		"lte":      "{0} must be at most {1}",
		"url":      "{0} must be a valid URL",
		"uuid":     "{0} must be a valid UUID",
	}

	for tag, msg := range messages {
		registerTranslation(tag, msg)
	}
}

func registerTranslation(tag, message string) {
	_ = validate.RegisterTranslation(tag, trans,
		func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field(), fe.Param())
			return t
		},
	)
}
