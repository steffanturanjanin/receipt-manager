package errors

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type HttpClientErrorInterface interface {
	AppErrorInterface
	GetCode() int
	GetErrors() map[string]string
}

type HttpError struct {
	ErrBase
	Code   int               `json:"code"`
	Errors map[string]string `json:"errors,omitempty"`
}

func (e HttpError) GetCode() int {
	return e.Code
}

func (e HttpError) GetErrors() map[string]string {
	return e.Errors
}

func NewHttpError(err error) *HttpError {
	appError, ok := err.(AppErrorInterface)
	if !ok {
		return nil
	}

	statusCode := HttpStatusCodeFromAppError(err)
	if statusCode == http.StatusInternalServerError {
		return nil
	}

	return &HttpError{
		ErrBase: ErrBase{
			Err:     appError.GetError(),
			Message: appError.GetMessage(),
		},
		Code: statusCode,
	}
}

func NewHttpValidationError(err error, trans ut.Translator) *HttpError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	validationErrorMap := map[string]string{}
	for _, e := range validationErrors {
		validationErrorMap[e.Field()] = e.Translate(trans)
	}

	return &HttpError{
		ErrBase: ErrBase{
			Err:     err,
			Message: "Input validation failed.",
		},
		Code:   http.StatusBadRequest,
		Errors: validationErrorMap,
	}
}

func HttpStatusCodeFromAppError(err error) int {
	switch err.(type) {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrResourceNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrDuplicateEntry:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
