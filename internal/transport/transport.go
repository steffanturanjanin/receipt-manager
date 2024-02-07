package transport

import "net/http"

type ValidationError struct {
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Errors  map[string]string `json:"errors"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func NewBadRequestResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
		Code:  http.StatusBadRequest,
	}
}

func NewValidationError(errors map[string]string) ValidationError {
	return ValidationError{
		Message: "Validation error.",
		Code:    http.StatusUnprocessableEntity,
		Errors:  errors,
	}
}
