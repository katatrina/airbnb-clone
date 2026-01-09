package validator

import (
	"reflect"
	"strings"
)

// NormalizeStruct preprocesses struct fields based on "normalize" tag.
// This function is framework-agnostic and can be used with any HTTP framework.
//
// Supported normalize rules:
//   - trim: Remove leading/trailing whitespace
//   - lower: Convert to lowercase
//   - upper: Convert to uppercase
//
// Example:
//
//	type User struct {
//	    Email string `normalize:"trim,lower"`
//	    Name  string `normalize:"trim"`
//	}
//
//	user := User{Email: "  TEST@Test.com  ", Name: "  John  "}
//	NormalizeStruct(&user)
//	// Result: user.Email = "test@test.com", user.Name = "John"
func NormalizeStruct(s interface{}) {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		normalizeTag := fieldType.Tag.Get("normalize")
		if normalizeTag == "" {
			continue
		}

		// Handle string fields
		if field.Kind() == reflect.String {
			str := field.String()

			// Parse normalize tag: "trim,lower"
			rules := strings.Split(normalizeTag, ",")
			for _, rule := range rules {
				switch strings.TrimSpace(rule) {
				case "trim":
					str = strings.TrimSpace(str)
				case "lower":
					str = strings.ToLower(str)
				case "upper":
					str = strings.ToUpper(str)
				}
			}

			field.SetString(str)
		}

		// Recursive for nested structs
		if field.Kind() == reflect.Struct {
			NormalizeStruct(field.Addr().Interface())
		}
	}
}
