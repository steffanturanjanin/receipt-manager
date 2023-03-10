package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	v "github.com/steffanturanjanin/receipt-manager/internal/validator"
)

var (
	encodedServerErrResp []byte = json.RawMessage(`{"message":"Internal server error."}`)
)

func ParseBody(destination interface{}, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	return nil
}

func JsonResponse(w http.ResponseWriter, payload interface{}, status int) {
	if payload == nil {
		w.WriteHeader(status)
		return
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		handleInternalServerError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err = w.Write(encoded); err != nil {
		panic(err)
	}
}

func JsonErrorResponse(w http.ResponseWriter, err error) {
	if httpClientError, ok := err.(errors.HttpClientErrorInterface); ok {
		JsonResponse(w, httpClientError, httpClientError.GetCode())

		return
	}

	handleInternalServerError(w, err)
}

func handleInternalServerError(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	if _, err := w.Write(encodedServerErrResp); err != nil {
		panic(err)
	}
}

// Valid validates the given struct.
func ValidateRequest(request interface{}, v *v.Validator) error {
	err := v.Struct(request)
	if err == nil {
		return nil
	}

	validationErrors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		validationErrors[fieldName] = err.Translate(v.GetTranslator())
	}

	return &errors.HttpError{
		ErrBase: errors.ErrBase{Err: err, Message: "Request validation failed."},
		Code:    http.StatusBadRequest,
		Errors:  validationErrors,
	}
}
