package query

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

type ReceiptFilters struct {
	Categories []int     `json:"categories"`
	Stores     []int     `json:"stores"`
	FromDate   time.Time `json:"from"`
	ToDate     time.Time `json:"to"`
	FromAmount *int      `json:"fromAmount"`
	ToAmount   *int      `json:"toAmount"`
}

const (
	// Filters
	CATEGORIES  = "categories"
	STORES      = "stores"
	FROM_DATE   = "fromDate"
	TO_DATE     = "toDate"
	FROM_AMOUNT = "fromAmount"
	TO_AMOUNT   = "toAmount"
)

type ReceiptQueryBuilder struct {
	BaseQueryBuilder
}

func NewReceiptQueryBuilder(query *gorm.DB) ReceiptQueryBuilder {
	query = query.Model(&models.Receipt{})
	baseQueryBuilder := BaseQueryBuilder{Query: query}

	return ReceiptQueryBuilder{
		BaseQueryBuilder: baseQueryBuilder,
	}
}

func (qb ReceiptQueryBuilder) Sort(sortQuery *SortQuery) ReceiptQueryBuilder {
	qb.BaseQueryBuilder.Sort(sortQuery)

	return qb
}

func (qb ReceiptQueryBuilder) SortableOptions() SortableOptions {
	return []string{"total_purchase_amount", "date"}
}

func (qb ReceiptQueryBuilder) Filter(rf ReceiptFilters) ReceiptQueryBuilder {
	// Date range filter
	qb.Query.Where("date BETWEEN ? AND ?", rf.FromDate, rf.ToDate)

	// Categories filter
	if rf.Categories != nil {
		qb.Query.Preload("ReceiptItems").Preload("ReceiptItems.Category").Where("category.id IN (?)", rf.Categories)
	}
	// Stores filter
	if len(rf.Stores) > 0 {
		qb.Query.Where("store_id IN (?)", rf.Stores)
	}
	// From amount filter
	if rf.FromAmount != nil {
		qb.Query.Where("total_purchase_amount > ?", *rf.FromAmount)
	}
	// To amount filter
	if rf.ToAmount != nil {
		qb.Query.Where("total_purchase_amount < ?", *rf.ToAmount)
	}

	return qb
}

func (qb ReceiptQueryBuilder) GetFilters(r *http.Request) ReceiptFilters {
	queryParams := r.URL.Query()

	// Extract categories filter
	var categories []int
	if categoriesParam := queryParams[CATEGORIES]; len(categoriesParam) > 0 {
		for _, categoryStr := range categoriesParam {
			if categoryInt, err := strconv.Atoi(categoryStr); err == nil {
				categories = append(categories, categoryInt)
			}
		}
	}

	// Extract stores filter
	var stores []int
	if storesParam := queryParams[STORES]; len(storesParam) > 0 {
		for _, storeStr := range storesParam {
			if storeInt, err := strconv.Atoi(storeStr); err == nil {
				stores = append(stores, storeInt)
			}
		}
	}

	// Time layout constant
	timeFormat := "2006-01-02"

	// Extract `fromDate` date range filter
	var fromDate time.Time
	if parsedTime, err := time.Parse(timeFormat, queryParams.Get(FROM_DATE)); err == nil {
		// Set parsed date
		fromDate = parsedTime
	} else {
		fmt.Printf("ERROR_PARSING_FROM_DATE: %s\n", err)
		// Get the current date
		currentDate := time.Now()

		// Get the first date of the current month
		firstDateOfMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())

		// Set first date of the current month as default `fromDate` filter value
		fromDate = firstDateOfMonth
	}

	// Extract `toDate` date range filter
	var toDate time.Time
	if parsedTime, err := time.Parse(timeFormat, queryParams.Get(TO_DATE)); err == nil {
		// Set parsed date
		toDate = parsedTime
	} else {
		fmt.Printf("ERROR_PARSING_TO_DATE: %s\n", err)

		// Get the current date
		currentDate := time.Now()

		// Get the first date of the next month
		firstDateOfNextMonth := time.Date(currentDate.Year(), currentDate.Month()+1, 1, 0, 0, 0, 0, currentDate.Location())

		// Subtract one day from the first date of the next month to get the last date of the current month
		lastDateOfMonth := firstDateOfNextMonth.AddDate(0, 0, -1)

		// Set first date of the current month as default `fromDate` filter value
		toDate = lastDateOfMonth
	}

	// Extract `fromAmount` range filter
	var fromAmount *int
	if fromAmountParam := queryParams.Get(FROM_AMOUNT); fromAmountParam != "" {
		if fromAmountValue, err := strconv.Atoi(fromAmountParam); err == nil {
			*fromAmount = fromAmountValue
		}
	}

	// Extract `toAmount` range filter
	var toAmount *int
	if toAmountParam := queryParams.Get(TO_AMOUNT); toAmountParam != "" {
		if toAmountValue, err := strconv.Atoi(toAmountParam); err == nil {
			*toAmount = toAmountValue
		}
	}

	return ReceiptFilters{
		Categories: categories,
		Stores:     stores,
		FromDate:   fromDate,
		ToDate:     toDate,
		FromAmount: fromAmount,
		ToAmount:   toAmount,
	}
}

func (qb ReceiptQueryBuilder) GetTotalPurchaseAmount() (*int, error) {
	var total *int
	if dbErr := qb.cloneQuery().Select("IFNULL(SUM(total_purchase_amount), 0)").Scan(&total).Error; dbErr != nil {
		return nil, dbErr
	}

	return total, nil
}
