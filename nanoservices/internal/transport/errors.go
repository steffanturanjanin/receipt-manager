package transport

import "net/http"

// ERRORS
type ValidationError struct {
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Errors  map[string]string `json:"errors"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// Bad Request - 400
func NewBadRequestResponse(message string) ErrorResponse {
	return ErrorResponse{
		Error: message,
		Code:  http.StatusBadRequest,
	}
}

// Forbidden - 403
func NewForbiddenError() ErrorResponse {
	return ErrorResponse{
		Error: "Forbidden",
		Code:  http.StatusForbidden,
	}
}

// Not Found - 404
func NewNotFoundError() ErrorResponse {
	return ErrorResponse{
		Error: "Not found",
		Code:  http.StatusNotFound,
	}
}

// Unprocessable Entity - 422
func NewValidationError(errors map[string]string) ValidationError {
	return ValidationError{
		Message: "Validation error.",
		Code:    http.StatusUnprocessableEntity,
		Errors:  errors,
	}
}

// Service Unavailable - 503
func NewServiceUnavailableError() ErrorResponse {
	return ErrorResponse{
		Error: "Service unavailable",
		Code:  http.StatusServiceUnavailable,
	}
}
