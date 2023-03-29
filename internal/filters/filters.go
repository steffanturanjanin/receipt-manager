package filters

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/steffanturanjanin/receipt-manager/internal/utils"
	"gorm.io/gorm"
)

type Operation string

const (
	FILTER_RANGE = "filter_range"
	FILTER_MATCH = "filter_match"
	FILTER_IN    = "filter_in"
	SEARCH       = "search"

	GreaterThan          Operation = ">"
	LessThan             Operation = "<"
	EqualTo              Operation = "="
	NotEqualTo           Operation = "<>"
	GreaterThanOrEqualTo Operation = ">="
	LessThanOrEqualTo    Operation = "<="
)

var FilterMatchOpetaionMapper = map[string]Operation{
	"gt":  GreaterThan,
	"lt":  LessThan,
	"eq":  EqualTo,
	"neq": NotEqualTo,
	"gte": GreaterThanOrEqualTo,
	"lte": LessThanOrEqualTo,
}

type FilterInterface interface {
	Filter(query *gorm.DB) *gorm.DB
}

type FiltersList []FilterInterface

type FilterMatch struct {
	Field     string
	Value     string
	Operation Operation
}

func (f FilterMatch) Filter(query *gorm.DB) *gorm.DB {
	queryString := fmt.Sprintf("%s %s ?", f.Field, f.Operation)
	query.Where(queryString, f.Value)

	return query
}

type FilterRange struct {
	Field     string
	LimitFrom string
	LimitTo   string
}

func (f FilterRange) Filter(query *gorm.DB) *gorm.DB {
	queryString := fmt.Sprintf("%s BETWEEN ? AND ?", f.Field)
	query.Where(queryString, f.LimitFrom, f.LimitTo)

	return query
}

type FilterIn struct {
	Field  string
	Values []string
}

func (f FilterIn) Filter(query *gorm.DB) *gorm.DB {
	queryString := fmt.Sprintf("%s IN ?", f.Field)
	return query.Where(queryString, f.Values)
}

type FilterableInterface interface {
	BuildFromRequest(r *http.Request) FilterableInterface
}

type abstractFilterable struct {
	FilterableInterface
	FiltersList FiltersList
}

func (f abstractFilterable) ApplyFilters(query *gorm.DB) *gorm.DB {
	for _, filters := range f.FiltersList {
		query = filters.Filter(query)
	}

	return query
}

func (f abstractFilterable) getFiltersRangeFromRequest(r *http.Request, allowedFields []string) []FilterInterface {
	filtersList := make([]FilterInterface, 0)

	if filtersRange, ok := r.URL.Query()[FILTER_RANGE]; ok {
		for _, filterRange := range filtersRange {
			if fr := getFilterRangeFromRequest(filterRange, allowedFields); fr != nil {
				filtersList = append(filtersList, *fr)
			}
		}
	}

	return filtersList
}

func getFilterRangeFromRequest(filterRange string, allowedFields []string) *FilterRange {
	params := strings.Split(filterRange, " ")
	field, value := params[0], params[1]

	if !utils.InSlice(allowedFields, field) {
		return nil
	}

	limits := strings.Split(value, ",")
	limitFrom, limitTo := limits[0], limits[1]

	return &FilterRange{
		Field:     field,
		LimitFrom: limitFrom,
		LimitTo:   limitTo,
	}
}

func (f abstractFilterable) getFiltersMatchFromRequest(r *http.Request, allowedFields []string) []FilterInterface {
	filtersList := make([]FilterInterface, 0)

	if filtersMatch, ok := r.URL.Query()[FILTER_MATCH]; ok {
		for _, filterMatch := range filtersMatch {
			if fm := getFilterMatchFromRequest(filterMatch, allowedFields); fm != nil {
				filtersList = append(filtersList, *fm)
			}
		}
	}

	return filtersList
}

func getFilterMatchFromRequest(filterMatch string, allowedFields []string) *FilterMatch {
	params := strings.Split(filterMatch, " ")
	field, comparator, value := params[0], params[1], params[2]

	if !utils.InSlice(allowedFields, field) {
		return nil
	}

	if operation, ok := FilterMatchOpetaionMapper[comparator]; ok {
		return &FilterMatch{
			Field:     field,
			Value:     value,
			Operation: operation,
		}
	}

	return nil
}

func (f abstractFilterable) getFiltersInFromRequest(r *http.Request, allowedFields []string) []FilterInterface {
	filterList := make([]FilterInterface, 0)

	if filtersIn, ok := r.URL.Query()[FILTER_IN]; ok {
		for _, filterIn := range filtersIn {
			if fi := getFilterInFromRequest(filterIn, allowedFields); fi != nil {
				filterList = append(filterList, *fi)
			}
		}
	}

	return filterList
}

func getFilterInFromRequest(filterMatch string, allowedFields []string) *FilterIn {
	params := strings.Split(filterMatch, " ")
	field, value := params[0], params[1]

	if !utils.InSlice(allowedFields, field) {
		return nil
	}

	if values := strings.Split(value, ","); len(values) > 0 {
		return &FilterIn{
			Field:  field,
			Values: values,
		}
	}

	return nil
}
