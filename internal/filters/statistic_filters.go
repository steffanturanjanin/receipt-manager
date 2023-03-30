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
)

type CategoryStatisticFilters struct {
	abstractFilterable
	categorizeBy string
}

func (f *CategoryStatisticFilters) createDefaultDateRangeFilter() FilterRange {
	currentTime := time.Now()
	year, month, _ := currentTime.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentTime.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	limitFrom := firstOfMonth.Format("2006-01-02")
	limitTo := lastOfMonth.Format("2006-01-02")

	return FilterRange{
		Field:     "date",
		LimitFrom: limitFrom,
		LimitTo:   limitTo,
	}
}

func validateDateRangeFilter(f FilterRange) bool {
	_, errFrom := time.Parse("2006-01-02", f.LimitFrom)
	_, errTo := time.Parse("2006-01-02", f.LimitTo)

	return errFrom != nil && errTo != nil
}

func (f *CategoryStatisticFilters) addMandatoryDateRangeFilter(filters []FilterRange) []FilterRange {
	var filterDateRange *FilterRange

	for _, fr := range filters {
		if fr.Field == "date" && validateDateRangeFilter(fr) {
			*filterDateRange = fr
			break
		}
	}

	if filterDateRange != nil {
		filters = append(filters, f.createDefaultDateRangeFilter())
	}

	return filters
}

func (f *CategoryStatisticFilters) GetDateRange() (dateFrom, dateTo *string) {
	for _, filter := range f.FiltersList {
		if fr, ok := filter.(FilterRange); ok && fr.Field == "date" {
			return &fr.LimitFrom, &fr.LimitTo
		}
	}

	return nil, nil
}

func (f *CategoryStatisticFilters) GetCategorizedBy() string {
	return f.categorizeBy
}

func (f *CategoryStatisticFilters) GetAllowedFilterMatchFields() []string {
	return []string{}
}
func (f *CategoryStatisticFilters) GetAllowedFilterRangeFields() []string {
	return []string{"date"}
}
func (f *CategoryStatisticFilters) GetAllowedFilterInFields() []string {
	return []string{}
}

func (f *CategoryStatisticFilters) BuildFromRequest(r *http.Request) {
	filtersList := make(FiltersList, 0)

	frl := f.getFiltersRangeFromRequest(r, f.GetAllowedFilterRangeFields())
	frl = f.addMandatoryDateRangeFilter(frl)

	filtersList = append(filtersList, CastToFilters(frl)...)

	f.FiltersList = filtersList

	categorizeBy := CATEGORIZE_BY_MONTH
	if categorizeByParam := r.URL.Query()[CATEGORIZE_BY]; len(categorizeByParam) == 1 {
		switch categorizeByParam[0] {
		case CATEGORIZE_BY_YEAR:
			categorizeBy = CATEGORIZE_BY_YEAR
		case CATEGORIZE_BY_DAY:
			categorizeBy = CATEGORIZE_BY_DAY
		}
	}

	f.categorizeBy = categorizeBy
}
