package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/katatrina/airbnb-clone/pkg/response"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

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

func TranslateErrors(err error) []response.FieldError {
	var fieldErrors []response.FieldError

	var validatorErrors validator.ValidationErrors
	if !errors.As(err, &validatorErrors) {
		return []response.FieldError{
			{
				Field:   "body",
				Code:    "INVALID_JSON",
				Message: "Invalid JSON format",
			},
		}
	}

	for _, e := range validatorErrors {
		fieldErrors = append(fieldErrors, response.FieldError{
			Field:   e.Field(),
			Code:    strings.ToUpper(e.Tag()),
			Message: e.Translate(trans),
		})
	}

	return fieldErrors
}
