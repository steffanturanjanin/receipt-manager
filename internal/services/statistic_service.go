package services

import (
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/filters"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
)

type StatisticService struct {
	statisticRepository repositories.StatisticRepositoryInterface
	categoryService     *CategoryService
}

func NewStatisticService(sr repositories.StatisticRepositoryInterface, cs *CategoryService) *StatisticService {
	return &StatisticService{
		statisticRepository: sr,
		categoryService:     cs,
	}
}

func (s *StatisticService) GetCategoryStatistic(f filters.CategoryStatisticFilters) (*CategoryStatistics, error) {
	categoryStatisticMap, err := s.statisticRepository.GetCategoryStatistic(f)
	if err != nil {
		return nil, err
	}

	categories, err := s.categoryService.GetAll()
	if err != nil {
		return nil, err
	}

	categoryStatistics := s.createCategoryStatistics(categories, f, categoryStatisticMap)

	return &categoryStatistics, nil
}

const (
	DAY_CATEGORIZATION_FORMAT   = "2006-01-02"
	MONTH_CATEGORIZATION_FORMAT = "2006-01"
	YEAR_CATEGORIZATION_FORMAT  = "2006"
)

func (s *StatisticService) createCategoryStatistics(c []dto.Category, f filters.CategoryStatisticFilters, m map[string]map[string]int) CategoryStatistics {
	dateFromString, dateToString := f.GetDateRange()

	dateFrom, _ := time.Parse("2006-01-02", *dateFromString)
	dateTo, _ := time.Parse("2006-01-02", *dateToString)

	categoryStatistics := CategoryStatistics{
		CategoryStatistics: make([]CategoryStatistic, 0),
	}

	categoryDateStatistic := newCategoryDateStatistic(f.GetCategorizedBy())

	for _, category := range c {
		categoryName := category.Name
		categoryStatistic := CategoryStatistic{
			CategoryName: categoryName,
			Stats:        make([]DateTotalAmount, 0),
		}

		for date := dateFrom; !date.After(dateTo); date = categoryDateStatistic.Next(date) {
			dateString := categoryDateStatistic.Format(date)
			dateTotalAmount := DateTotalAmount{
				Date:   dateString,
				Amount: m[categoryName][dateString],
			}

			categoryStatistic.Stats = append(categoryStatistic.Stats, dateTotalAmount)
		}

		categoryStatistics.CategoryStatistics = append(categoryStatistics.CategoryStatistics, categoryStatistic)
	}

	return categoryStatistics
}

func newCategoryDateStatistic(categorizeBy string) categoryDateStatisticInterface {
	switch categorizeBy {
	case filters.CATEGORIZE_BY_DAY:
		return categoryDailyStatistic{}
	case filters.CATEGORIZE_BY_MONTH:
		return categoryMonthlyStatistic{}
	case filters.CATEGORIZE_BY_YEAR:
		return categoryYearlyStatistic{}
	default:
		return categoryMonthlyStatistic{}
	}
}

type categoryDateStatisticInterface interface {
	Next(time.Time) time.Time
	Format(time.Time) string
}

type categoryDailyStatistic struct{}

func (c categoryDailyStatistic) Next(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

func (c categoryDailyStatistic) Format(t time.Time) string {
	return t.Format("2006-01-02")
}

type categoryMonthlyStatistic struct{}

func (c categoryMonthlyStatistic) Next(t time.Time) time.Time {
	return t.AddDate(0, 1, 0)
}

func (c categoryMonthlyStatistic) Format(t time.Time) string {
	return t.Format("2006-01")
}

type categoryYearlyStatistic struct{}

func (c categoryYearlyStatistic) Next(t time.Time) time.Time {
	return t.AddDate(1, 0, 0)
}

func (c categoryYearlyStatistic) Format(t time.Time) string {
	return t.Format("2006")
}

type DateTotalAmount struct {
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

type CategoryStatistic struct {
	CategoryName string            `json:"category_name"`
	Stats        []DateTotalAmount `json:"statistic"`
}

type CategoryStatistics struct {
	CategoryStatistics []CategoryStatistic `json:"categories"`
}
