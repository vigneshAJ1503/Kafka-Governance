package utils

import "net/http"

// ErrorType defines the type of error
type ErrorType int

const (
	ErrInvalidInput ErrorType = iota
	ErrNotFound
	ErrAlreadyExists
	ErrUnauthorized
	ErrForbidden
	ErrInternalServer
)

// APIError represents an API error with status code
type APIError struct {
	Type       ErrorType
	Message    string
	StatusCode int
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewInvalidInputError creates a 400 Bad Request error
func NewInvalidInputError(message string) *APIError {
	return &APIError{
		Type:       ErrInvalidInput,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a 404 Not Found error
func NewNotFoundError(message string) *APIError {
	return &APIError{
		Type:       ErrNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewAlreadyExistsError creates a 409 Conflict error
func NewAlreadyExistsError(message string) *APIError {
	return &APIError{
		Type:       ErrAlreadyExists,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewUnauthorizedError creates a 401 Unauthorized error
func NewUnauthorizedError(message string) *APIError {
	return &APIError{
		Type:       ErrUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a 403 Forbidden error
func NewForbiddenError(message string) *APIError {
	return &APIError{
		Type:       ErrForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewInternalServerError creates a 500 Internal Server Error
func NewInternalServerError(message string) *APIError {
	return &APIError{
		Type:       ErrInternalServer,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) (*APIError, bool) {
	apiErr, ok := err.(*APIError)
	return apiErr, ok
}
