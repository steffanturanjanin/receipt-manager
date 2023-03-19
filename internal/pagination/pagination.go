package pagination

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/steffanturanjanin/receipt-manager/internal/utils"
)

const (
	Ascending  string = "asc"
	Descending string = "desc"

	MaxLimit     int    = 100
	DefaultLimit int    = 25
	DefaultPage  int    = 1
	DefaultSort  string = Ascending

	LimitQueryParam string = "limit"
	PageQueryParam  string = "page"
	SortQueryParam  string = "sort"
)

func isLimitValueValid(limit int) bool {
	return limit > 0 && limit <= MaxLimit
}

type Pagination struct {
	Limit       int
	Page        int
	Sort        SortList
	TotalEnties int
	TotalPages  int
}

type SortList []Sort

func ValidateSortList(sl SortList, s Sortable) {
	for sortListElIndex, sortListEl := range sl {
		if utils.InSlice(s.SortableFields(), sortListEl.Field) {
			continue
		}

		sl = append(sl[:sortListElIndex], sl[sortListElIndex+1:]...)
	}
}

type Sort struct {
	Field     string
	Direction string
}

type Sortable interface {
	SortableFields() []string
}

func GetPaginationFromRequest(r *http.Request) Pagination {
	limit := DefaultLimit
	page := DefaultPage
	sortList := make(SortList, 0)

	queryParams := r.URL.Query()

	if limitParam := queryParams.Get(LimitQueryParam); limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err == nil && isLimitValueValid(l) {
			limit = l
		}
	}

	if pageParam := queryParams.Get(PageQueryParam); pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err == nil && page > 0 {
			page = p
		}
	}

	if sort := queryParams.Get(SortQueryParam); sort != "" {
		sortParams := strings.Split(sort, ",")
		for _, sortParam := range sortParams {
			sortParamPair := strings.Split(sortParam, " ")
			field := sortParamPair[0]
			direction := sortParamPair[1]

			// if !contains(sortable.SortableFields(), field) {
			// 	continue
			// }

			if direction != Ascending && direction != Descending {
				direction = DefaultSort
			}

			sortList = append(sortList, Sort{Field: field, Direction: direction})
		}
	}

	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sortList,
	}
}
