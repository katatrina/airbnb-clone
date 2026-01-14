package request

import (
	"reflect"
	"strings"
)

// NormalizeStruct applies normalization rules from struct field tags.
// Supported rules: trim, lower, upper, singlespace.
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
		if !field.CanSet() {
			continue
		}

		normalizeTag := t.Field(i).Tag.Get("normalize")
		if normalizeTag == "" || field.Kind() != reflect.String {
			continue
		}

		str := field.String()
		rules := strings.Split(normalizeTag, ",")

		for _, rule := range rules {
			switch strings.TrimSpace(rule) {
			case "trim":
				str = strings.TrimSpace(str)
			case "lower":
				str = strings.ToLower(str)
			case "upper":
				str = strings.ToUpper(str)
			case "singlespace":
				str = strings.Join(strings.Fields(str), " ")
			}
		}

		field.SetString(str)
	}
}
