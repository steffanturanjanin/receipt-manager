package repositories

import (
	"fmt"

	"github.com/steffanturanjanin/receipt-manager/internal/filters"
	"gorm.io/gorm"
)

type StatisticRepositoryInterface interface {
	GetCategoryStatistic(filters.CategoryStatisticFilters) (map[string]map[string]int, error)
}

type StatisticRepository struct {
	db *gorm.DB
}

func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{
		db: db,
	}
}

func (r *StatisticRepository) GetCategoryStatistic(f filters.CategoryStatisticFilters) (map[string]map[string]int, error) {
	categoryStatistics := make(map[string]map[string]int)

	baseQuery := r.db.Table("categories as c").
		Select("c.name as name", "SUM(ri.total_amount) as total", fmt.Sprintf("%s as date", r.formDateString(f, "r.date"))).
		Joins("LEFT JOIN receipt_items as ri ON c.id = ri.category_id").
		Joins("LEFT JOIN receipts as r ON ri.receipt_id = r.id").
		Group("c.id, date")

	rows, err := f.ApplyFilters(baseQuery).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var categoryName string
		var total int
		var receiptDate string

		if err := rows.Scan(&categoryName, &total, &receiptDate); err != nil {
			continue
		}

		if _, ok := categoryStatistics[categoryName]; !ok {
			categoryStatistics[categoryName] = make(map[string]int, 0)
		}

		categoryStatistics[categoryName][receiptDate] = total
	}

	return categoryStatistics, nil
}

func (r *StatisticRepository) formDateString(f filters.CategoryStatisticFilters, field string) string {
	var date string

	switch f.GetCategorizedBy() {
	case filters.CATEGORIZE_BY_YEAR:
		date = fmt.Sprintf("CONCAT_WS('-', YEAR(%[1]v))", field)
	case filters.CATEGORIZE_BY_MONTH:
		date = fmt.Sprintf("CONCAT_WS('-', YEAR(%[1]v), LPAD(MONTH(%[1]v), 2, '0'))", field)
	case filters.CATEGORIZE_BY_DAY:
		date = fmt.Sprintf("CONCAT_WS('-', YEAR(%[1]v), LPAD(MONTH(%[1]v), 2, '0'), LPAD(DAY(%[1]v), 2, '0'))", field)
	}

	return date
}
