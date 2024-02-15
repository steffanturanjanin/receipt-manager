package transport

import (
	"net/http"
)

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
func NewBadRequestResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
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

// Unprocessable entity - 422
func NewValidationError(errors map[string]string) ValidationError {
	return ValidationError{
		Message: "Validation error.",
		Code:    http.StatusUnprocessableEntity,
		Errors:  errors,
	}
}

// Responses
type PaginationResponse struct {
	Data []interface{} `json:"data"`
	Meta interface{}   `json:"meta"`
}

func TransformResponseData[T any](models []T) []interface{} {
	var data []interface{}
	for _, model := range models {
		data = append(data, model)
	}

	return data
}

func CreatePaginationResponse(data []interface{}, meta interface{}) (*PaginationResponse, error) {
	return &PaginationResponse{Data: data, Meta: meta}, nil
}
