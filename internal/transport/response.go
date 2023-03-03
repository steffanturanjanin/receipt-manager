package transport

import (
	"encoding/json"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(data)
}

func ValidationErrorResponseJson(w http.ResponseWriter, err error) {
	validationErrorResponse := NewValidationErrorResponse(err)

	ResponseJson(w, validationErrorResponse, validationErrorResponse.Code)
}
