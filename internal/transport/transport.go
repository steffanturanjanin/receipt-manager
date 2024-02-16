package transport

import (
	"errors"
	"reflect"
)

// Pagination
type PaginationResponse struct {
	Data []interface{} `json:"data"`
	Meta interface{}   `json:"meta"`
}

// Builds pagination response. Data param needs to be a slice.
// Converts interface{} param to []interface{} expected by `PaginationResponse` struct.
func CreatePaginationResponse(data interface{}, meta interface{}) (*PaginationResponse, error) {
	// Check if destination is a pointer pointing to a slice
	value := reflect.ValueOf(data).Elem()
	if value.Kind() != reflect.Slice {
		return nil, errors.New("data must point to a slice")
	}

	items := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		items[i] = value.Index(i).Interface()
	}

	return &PaginationResponse{Data: items, Meta: meta}, nil
}
