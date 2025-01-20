package helper

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, err := range errs {
		errors[err.Field()] = err.Tag()
	}
	return errors
}

func ValidateStruct(validate *validator.Validate, log *logrus.Logger, request interface{}) error {
	if err := validate.Struct(request); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			log.Warnf("Validation failed: %+v", validationErrors)
			formattedErrors := FormatValidationErrors(validationErrors)
			return model.ErrValidationFailed(formattedErrors)
		}
		log.Warnf("Failed to validate request body: %+v", err)
		return model.ErrBadRequest
	}
	return nil
}
