package filters

import (
	"net/http"
)

type ReceiptFilters struct {
	abstractFilterable
}

func (f *ReceiptFilters) GetAllowedFilterRangeFields() []string {
	return []string{
		"total_purchase_amount",
		"date",
	}
}

func (f *ReceiptFilters) GetAllowedFilterMatchFields() []string {
	return []string{
		"total_purchase_amount",
		"date",
	}
}

func (f *ReceiptFilters) GetAllowedFilterInFields() []string {
	return []string{
		"date",
		"store_id",
	}
}

func (f *ReceiptFilters) BuildFromRequest(r *http.Request) {
	filtersList := make(FiltersList, 0)

	frl := f.getFiltersRangeFromRequest(r, f.GetAllowedFilterRangeFields())
	fml := f.getFiltersMatchFromRequest(r, f.GetAllowedFilterMatchFields())
	fil := f.getFiltersInFromRequest(r, f.GetAllowedFilterInFields())

	filtersList = append(filtersList, CastToFilters(frl)...)
	filtersList = append(filtersList, CastToFilters(fml)...)
	filtersList = append(filtersList, CastToFilters(fil)...)

	f.FiltersList = filtersList
}
