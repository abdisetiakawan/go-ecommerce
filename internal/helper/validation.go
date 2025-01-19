package helper

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, err := range errs {
		errors[err.Field()] = err.Tag()
	}
	return errors
}