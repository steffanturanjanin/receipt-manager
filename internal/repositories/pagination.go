package repositories

import (
	"fmt"
	"math"

	"github.com/steffanturanjanin/receipt-manager/internal/pagination"
	"gorm.io/gorm"
)

func Paginate(query *gorm.DB, p *pagination.Pagination, s pagination.Sortable) *gorm.DB {
	pagination.ValidateSortList(p.Sort, s)

	var totalRecords int64
	var totalPages int

	query.Count(&totalRecords)

	totalPages = int(math.Ceil(float64(totalRecords) / float64(p.Limit)))

	p.TotalEntries = int(totalRecords)
	p.TotalPages = totalPages

	offset := (p.Page - 1) * p.Limit

	sortString := ""
	for index, sortElement := range p.Sort {
		sortString = sortString + fmt.Sprintf("%s %s", sortElement.Field, sortElement.Direction)
		if index != len(p.Sort)-1 {
			sortString = sortString + ", "
		}
	}

	query.Offset(offset).Limit(p.Limit).Order(sortString)

	return query
}
