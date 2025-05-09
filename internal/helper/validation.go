package helper

import (
	"fmt"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors  map[string]ValidationError `json:"errors"`
	Message string                     `json:"message"`
	Status  string                     `json:"status"`
}

func GetValidationMessage(field string, tag string, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, param)
	case "oneof":
		if field == "PaymentMethod" {
			return "Payment method must be either 'cash' or 'transfer'"
		}
		if field == "Gender" {
			return "Gender must be either 'male' or 'female'"
		}
		if field == "Category" {
			return "Category must be either 'men' or 'women'"
		}
		return fmt.Sprintf("%s has invalid value", field)
	case "e164":
		return "Phone number must be in E.164 format"
	case "url":
		return "Must be a valid URL"
	case "numeric":
		return "Must be a number"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", field, param)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func FormatValidationErrors(errs validator.ValidationErrors) *ValidationErrors {
	errors := make(map[string]ValidationError)

	for _, err := range errs {
		errors[err.Field()] = ValidationError{
			Tag:     err.Tag(),
			Message: GetValidationMessage(err.Field(), err.Tag(), err.Param()),
		}
	}

	return &ValidationErrors{
		Errors:  errors,
		Message: "Validation failed",
		Status:  "fail",
	}
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
