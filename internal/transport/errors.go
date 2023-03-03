package transport

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []fieldError `json:"errors"`
	HttpErrorResponse
}

type HttpErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewValidationErrorResponse(err error) ValidationErrorResponse {
	validationError := ValidationErrorResponse{
		HttpErrorResponse: HttpErrorResponse{
			Message: "Validation failed.",
			Code:    http.StatusBadRequest,
		},
	}

	for _, err := range err.(validator.ValidationErrors) {
		fieldError := fieldError{
			Field:   err.Field(),
			Message: errorMessageForValidationTag(err.Tag()),
		}

		validationError.Errors = append(validationError.Errors, fieldError)
	}

	return validationError
}

func errorMessageForValidationTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required."
	case "email":
		return "Invalid email."
	default:
		return ""
	}
}
