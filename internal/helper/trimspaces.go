package helper

import (
	"reflect"
	"regexp"
	"strings"
)

func TrimSpaces(s interface{}, excludeValues ...string) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}

	excludeMap := make(map[string]bool)
	for _, value := range excludeValues {
		excludeMap[value] = true
	}

	space := regexp.MustCompile(`\s+`)

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			currentValue := field.String()
			if !excludeMap[currentValue] {
				normalized := space.ReplaceAllString(strings.TrimSpace(currentValue), " ")
				normalized = strings.ToLower(normalized)
				field.SetString(normalized)
			}
		}
	}
}