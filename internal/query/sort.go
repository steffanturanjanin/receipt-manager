package query

import (
	"net/http"

	"gorm.io/gorm"
)

type SortableOptions = []string

type SortDirection string

type Sortable interface {
	SortableOptions() SortableOptions
	Sort(*SortQuery) *gorm.DB
}

type SortQuery struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

const (
	ASCENDING    SortDirection = "ASC"
	DESCENDING   SortDirection = "DESC"
	DEFAULT_SORT SortDirection = ASCENDING
)

const (
	SORT_PARAM      = "sort"
	DIRECTION_PARAM = "direction"
)

func SortQueryFromRequest(r http.Request) *SortQuery {
	queryParams := r.URL.Query()
	sort := queryParams.Get(SORT_PARAM)
	direction := queryParams.Get(DIRECTION_PARAM)

	if sort == "" {
		return nil
	}

	if direction == "" {
		direction = string(DEFAULT_SORT)
	}

	return &SortQuery{Field: sort, Direction: SortDirection(direction)}
}
