package query

import (
	"math"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

const (
	// Page param
	PAGE_PARAM         = "page"
	DEFAULT_PAGE_PARAM = 1

	// Limit Param
	LIMIT_PARAM   = "limit"
	DEFAULT_LIMIT = 25
	MAX_LIMIT     = 100
)

type PaginationQuery struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginationMeta struct {
	Page         int  `json:"page"`
	PrevPage     *int `json:"prevPage"`
	NextPage     *int `json:"nextPage"`
	PerPage      int  `json:"perPage"`
	TotalPages   int  `json:"totalPages"`
	TotalEntries int  `json:"totalEntries"`
}

type PaginationData struct {
	Data []interface{}
	Meta PaginationMeta
}

func PaginationQueryFromRequest(r http.Request) PaginationQuery {
	queryParams := r.URL.Query()
	pageParam := queryParams.Get(PAGE_PARAM)
	limitParam := queryParams.Get(LIMIT_PARAM)

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || (limit < 0 || limit > MAX_LIMIT) {
		limit = DEFAULT_LIMIT
	}

	return PaginationQuery{Page: page, Limit: limit}
}

func GetPaginationMeta(query *gorm.DB, pagination PaginationQuery) PaginationMeta {
	var totalEntries int64
	var totalPages int
	var prevPage *int
	var nextPage *int

	// Total entries
	query.Count(&totalEntries)

	// Total pages
	totalPages = int(math.Ceil(float64(totalEntries) / float64(pagination.Limit)))

	// Prev page
	if pagination.Page > 1 {
		*prevPage = pagination.Page - 1
	}

	// Next page
	if pagination.Page < totalPages {
		*nextPage = pagination.Page + 1
	}

	return PaginationMeta{
		Page:         pagination.Page,
		PerPage:      pagination.Limit,
		PrevPage:     prevPage,
		NextPage:     nextPage,
		TotalPages:   totalPages,
		TotalEntries: int(totalEntries),
	}
}
