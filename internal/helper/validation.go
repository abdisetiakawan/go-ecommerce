package helper

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, err := range errs {
		errors[err.Field()] = err.Tag()
	}
	return errors
}

func ValidateStruct(validate *validator.Validate, request interface{}) error {
	if err := validate.Struct(request); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := FormatValidationErrors(validationErrors)
			return model.ErrValidationFailed(formattedErrors)
		}
		return model.ErrBadRequest
	}
	return nil
}
