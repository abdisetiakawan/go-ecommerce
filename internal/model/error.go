package model

import "github.com/gofiber/fiber/v2"

type ApiError struct {
    StatusCode int         `json:"-"`
    Message    string      `json:"message"`
    Errors     interface{} `json:"errors,omitempty"`
}

func (e *ApiError) Error() string {
    return e.Message
}

func NewApiError(statusCode int, message string, errors interface{}) *ApiError {
    return &ApiError{
        StatusCode: statusCode,
        Message:    message,
        Errors:     errors,
    }
}

var (
    ErrUserAlreadyExists  = NewApiError(fiber.StatusConflict, "User already exists", nil)
    ErrInvalidCredentials = NewApiError(fiber.StatusUnauthorized, "Invalid credentials", nil)
    ErrBadRequest         = NewApiError(fiber.StatusBadRequest, "Invalid request", nil)
    ErrInternalServer     = NewApiError(fiber.StatusInternalServerError, "Internal server error", nil)
    ErrNotFound           = NewApiError(fiber.StatusNotFound, "Resource not found", nil)
    ErrConflict           = NewApiError(fiber.StatusConflict, "Conflict", nil)
    ErrUsernameExists     = NewApiError(fiber.StatusConflict, "Username already exists", nil)
    ErrForbidden          = NewApiError(fiber.StatusForbidden, "You are not allowed to access this resource", nil)
    ErrPasswordNotMatch   = NewApiError(fiber.StatusBadRequest, "Password and confirm password do not match", nil)
    ErrStoreNotFound      = NewApiError(fiber.StatusNotFound, "Store not found", nil)
)

func ErrValidationFailed(errors interface{}) *ApiError {
    return NewApiError(fiber.StatusBadRequest, "Validation failed", errors)
}