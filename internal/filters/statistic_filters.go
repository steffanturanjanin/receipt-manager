package filters

import (
	"net/http"
	"time"
)

const (
	CATEGORIZE_BY       = "categorize_by"
	CATEGORIZE_BY_YEAR  = "year"
	CATEGORIZE_BY_MONTH = "month"
	CATEGORIZE_BY_DAY   = "day"

	DATE_FROM = "date_from"
	DATE_TO   = "date_to"
)

type dateRange struct {
	From string
	To   string
}

//type categoryIds []string

type CategoryStatisticFilters struct {
	FilterDateRange dateRange
	CategorizeBy    string
}

func (f *CategoryStatisticFilters) BuildFromRequest(r *http.Request) {
	currentTime := time.Now()
	year, month, _ := currentTime.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentTime.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	limitFrom := firstOfMonth.Format("2006-01-02")
	limitTo := lastOfMonth.Format("2006-01-02")

	dateFromParam := r.URL.Query()[DATE_FROM]
	dateToParam := r.URL.Query()[DATE_TO]

	if len(dateFromParam) == 1 && len(dateToParam) == 1 {
		_, errFrom := time.Parse("2006-01-02", dateFromParam[0])
		_, errTo := time.Parse("2006-01-02", dateToParam[0])

		if errFrom == nil && errTo == nil {
			limitFrom = dateFromParam[0]
			limitTo = dateToParam[0]
		}
	}

	f.FilterDateRange = dateRange{
		From: limitFrom,
		To:   limitTo,
	}

	categorizeBy := CATEGORIZE_BY_MONTH
	if categorizeByParam := r.URL.Query()[CATEGORIZE_BY]; len(categorizeByParam) == 1 {
		switch categorizeByParam[0] {
		case CATEGORIZE_BY_YEAR:
			categorizeBy = CATEGORIZE_BY_YEAR
		case CATEGORIZE_BY_DAY:
			categorizeBy = CATEGORIZE_BY_DAY
		}
	}

	f.CategorizeBy = categorizeBy
}

// func (f *CategoryStatisticFilters) BuildFromRequest(r *http.Request) {
// 	filtersList := make(FiltersList, 0)

// 	currentTime := time.Now()
// 	year, month, _ := currentTime.Date()
// 	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentTime.Location())
// 	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

// 	limitFrom := firstOfMonth.Format("2006-01-02")
// 	limitTo := lastOfMonth.Format("2006-01-02")

// 	for _, rangeFilter := range r.URL.Query()[FILTER_RANGE] {
// 		params := strings.Split(" ", rangeFilter)
// 		field, value := params[0], params[1]
// 		rangeValue := strings.Split(value, ",")

// 		if field == "date" {
// 			limitFrom, limitTo = rangeValue[0], rangeValue[1]
// 		}
// 	}

// 	filtersList = append(filtersList, FilterRange{
// 		Field:     "date",
// 		LimitFrom: limitFrom,
// 		LimitTo:   limitTo,
// 	})

// 	f.FiltersList = filtersList

// 	categorizeBy := CATEGORIZE_BY_MONTH
// 	if categorizeByParam := r.URL.Query()[CATEGORIZE_BY]; len(categorizeBy) == 1 {
// 		switch categorizeByParam[0] {
// 		case CATEGORIZE_BY_YEAR:
// 			categorizeBy = CATEGORIZE_BY_YEAR
// 		case CATEGORIZE_BY_DAY:
// 			categorizeBy = CATEGORIZE_BY_DAY
// 		}
// 	}

// 	f.CategorizeBy = categorizeBy
// }
