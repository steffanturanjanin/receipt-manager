package repositories

import (
	"fmt"

	"github.com/steffanturanjanin/receipt-manager/internal/filters"
	"gorm.io/gorm"
)

type StatisticRepositoryInterface interface {
	GetCategoryStatistic(filters.CategoryStatisticFilters) (map[string]map[string]int, error)
	GetStoreStatisticsForCategory(int, filters.StoreStatisticForCategoryFilters) ([]StoreStatisticItem, error)
	GetCategoryStatisticsForStore(string, filters.CategoryStatisticsForStoreFilters) ([]CategoryStatisticForStore, error)
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

type StoreStatisticItem struct {
	Store struct {
		Tin  string
		Name string
	}
	ReceiptItem struct {
		Name  string
		Price int
	}
}

func (r *StatisticRepository) GetStoreStatisticsForCategory(categoryId int, f filters.StoreStatisticForCategoryFilters) ([]StoreStatisticItem, error) {
	baseQuery := r.db.Table("categories as c").
		Select("s.tin as store_tin", "s.name as store_name", "ri.name as article_name", "ri.total_amount as article_price").
		Joins("INNER JOIN receipt_items as ri ON c.id = ri.category_id").
		Joins("INNER JOIN receipts as r ON ri.receipt_id = r.id").
		Joins("INNER JOIN stores as s ON r.store_id = s.tin").
		Where("c.id = ?", categoryId)

	rows, err := f.ApplyFilters(baseQuery).Rows()
	if err != nil {
		return nil, err
	}

	result := []StoreStatisticItem{}

	defer rows.Close()
	for rows.Next() {
		var storeTin string
		var storeName string
		var articleName string
		var articlePrice int

		if err := rows.Scan(&storeTin, &storeName, &articleName, &articlePrice); err != nil {
			continue
		}

		result = append(result, StoreStatisticItem{
			Store: struct {
				Tin  string
				Name string
			}{
				Tin:  storeTin,
				Name: storeName,
			},
			ReceiptItem: struct {
				Name  string
				Price int
			}{
				Name:  articleName,
				Price: articlePrice,
			},
		})
	}

	return result, nil
}

type CategoryStatisticForStore struct {
	CategoryId   int
	CategoryName string
	Total        int
}

func (r *StatisticRepository) GetCategoryStatisticsForStore(storeTin string, f filters.CategoryStatisticsForStoreFilters) ([]CategoryStatisticForStore, error) {
	baseQuery := r.db.Table("stores as s").
		Select("c.id as category_id", "c.name as category_name", "sum(ri.total_amount) as total").
		Joins("INNER JOIN receipts as r ON s.tin = r.store_id").
		Joins("INNER JOIN receipt_items as ri ON r.id = ri.receipt_id").
		Joins("INNER JOIN categories as c ON ri.category_id = c.id").
		Where("s.tin = ?", storeTin).
		Group("c.id")

	rows, err := f.ApplyFilters(baseQuery).Rows()
	if err != nil {
		return nil, err
	}

	result := []CategoryStatisticForStore{}

	defer rows.Close()
	for rows.Next() {
		var categoryId int
		var categoryName string
		var total int

		if err := rows.Scan(&categoryId, &categoryName, &total); err != nil {
			continue
		}

		result = append(result, CategoryStatisticForStore{
			CategoryId:   categoryId,
			CategoryName: categoryName,
			Total:        total,
		})
	}

	return result, nil
}
