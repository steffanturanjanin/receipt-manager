package filters

import (
	"net/http"
)

var (
	AllowedFilterRangeFields = []string{
		"total_purchase_amount",
		"date",
	}

	AllowedFilterMatchFields = []string{
		"total_purchase_amount",
		"date",
	}

	AllowedFilterInFields = []string{
		"date",
		"store_id",
	}
)

type ReceiptFilters struct {
	abstractResourceFilters
}

func (f *ReceiptFilters) BuildFromRequest(r *http.Request) {
	filtersList := make(FiltersList, 0)

	frl := f.getFiltersRangeFromRequest(r, AllowedFilterRangeFields)
	fml := f.getFiltersMatchFromRequest(r, AllowedFilterMatchFields)
	fil := f.getFiltersInFromRequest(r, AllowedFilterInFields)

	filtersList = append(filtersList, frl...)
	filtersList = append(filtersList, fml...)
	filtersList = append(filtersList, fil...)

	f.FiltersList = filtersList
}
