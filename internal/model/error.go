package model

import "github.com/gofiber/fiber"

type ApiError struct {
    StatusCode int    `json:"-"`
    Message    string `json:"message"`
}

func (e *ApiError) Error() string {
    return e.Message
}

func NewApiError(statusCode int, message string) *ApiError {
    return &ApiError{
        StatusCode: statusCode,
        Message:    message,
    }
}

var (
    ErrUserAlreadyExists  = NewApiError(fiber.StatusConflict, "User already exists")
    ErrInvalidCredentials = NewApiError(fiber.StatusUnauthorized, "Invalid credentials")
    ErrBadRequest        = NewApiError(fiber.StatusBadRequest, "Invalid request")
    ErrInternalServer    = NewApiError(fiber.StatusInternalServerError, "Internal server error")
    ErrNotFound          = NewApiError(fiber.StatusNotFound, "Resource not found")
    ErrConflict = NewApiError(fiber.StatusConflict, "Conflict")
    ErrUsernameExists = NewApiError(fiber.StatusConflict, "Username already exists")
    ErrForbidden = NewApiError(fiber.StatusForbidden, "You are not allowed to access this resource")
    ErrPasswordNotMatch = NewApiError(fiber.StatusBadRequest, "Password and confirm password do not match")
)