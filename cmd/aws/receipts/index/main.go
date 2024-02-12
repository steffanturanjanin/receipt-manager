package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
)

const (
	// Query param keys
	PAGE  = "page"
	LIMIT = "limit"
	// Filters
	CATEGORIES  = "categories"
	STORES      = "stores"
	FROM_DATE   = "fromDate"
	TO_DATE     = "toDate"
	FROM_AMOUNT = "fromAmount"
	TO_AMOUNT   = "toAmount"
)

// Sort
type SortDirection string

const (
	ASC  SortDirection = "ASC"
	DESC SortDirection = "DESC"
)

func GetAllowedSortFields() []string {
	return []string{"date", "total_purchase_amount"}
}

type Response struct {
	Data []interface{} `json:"data"`
	Meta interface{}   `json:"meta"`
}

type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

type PaginationMeta struct {
	Page         int  `json:"page"`
	PrevPage     *int `json:"prevPage"`
	NextPage     *int `json:"nextPage"`
	PerPage      int  `json:"perPage"`
	TotalPages   int  `json:"totalPages"`
	TotalEntries int  `json:"totalEntries"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Filters struct {
	Categories []int     `json:"categories"`
	Stores     []int     `json:"stores"`
	FromDate   time.Time `json:"from"`
	ToDate     time.Time `json:"to"`
	FromAmount *int      `json:"fromAmount"`
	ToAmount   *int      `json:"toAmount"`
}

type ReceiptsMeta struct {
	Pagination PaginationMeta `json:"pagination"`
	Total      int            `json:"total"`
}

func GetPagination(queryParams url.Values) Pagination {
	var page int = 1
	var limit int = 25

	if pageParam := queryParams.Get(PAGE); pageParam != "" {
		if pageInt, err := strconv.Atoi(pageParam); err == nil && pageInt > 0 {
			page = pageInt
		}
	}
	if limitParam := queryParams.Get(LIMIT); limitParam != "" {
		if limitInt, err := strconv.Atoi(limitParam); err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	return Pagination{Page: page, Limit: limit}
}

func GetSort(queryParams url.Values) *Sort {
	field := queryParams.Get("sort")
	direction := queryParams.Get("direction")

	allowedFields := GetAllowedSortFields()
	for _, allowedField := range allowedFields {
		if field == allowedField {
			return &Sort{Field: field, Direction: SortDirection(direction)}
		}
	}

	return nil
}

func GetPaginationMeta(query *gorm.DB, pagination Pagination) PaginationMeta {
	var totalEntries int64
	var totalPages int
	var prevPage *int
	var nextPage *int

	// Total entries
	query.Count(&totalEntries)

	// Total pages
	totalPages = int(math.Ceil(float64(totalEntries) / float64(pagination.Limit)))

	// Prev page
	if pagination.Page > 1 {
		*prevPage = pagination.Page - 1
	}

	// Next page
	if pagination.Page < totalPages {
		*nextPage = pagination.Page + 1
	}

	return PaginationMeta{
		Page:         pagination.Page,
		PerPage:      pagination.Limit,
		PrevPage:     prevPage,
		NextPage:     nextPage,
		TotalPages:   totalPages,
		TotalEntries: int(totalEntries),
	}
}

func GetFilters(queryParams url.Values) Filters {
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
	timeFormat := "2006-01-02 15:04:05"

	// Extract `fromDate` date range filter
	var fromDate time.Time
	if parsedTime, err := time.Parse(queryParams.Get(FROM_DATE), timeFormat); err == nil {
		// Set parsed date
		fromDate = parsedTime
	} else {
		// Get the current date
		currentDate := time.Now()

		// Get the first date of the current month
		firstDateOfMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())

		// Set first date of the current month as default `fromDate` filter value
		fromDate = firstDateOfMonth
	}

	// Extract `toDate` date range filter
	var toDate time.Time
	if parsedTime, err := time.Parse(queryParams.Get(TO_DATE), timeFormat); err == nil {
		// Set parsed date
		toDate = parsedTime
	} else {
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

	return Filters{
		Categories: categories,
		Stores:     stores,
		FromDate:   fromDate,
		ToDate:     toDate,
		FromAmount: fromAmount,
		ToAmount:   toAmount,
	}
}

func BuildQuery(baseQuery *gorm.DB, pagination Pagination, filters Filters, sort *Sort) *gorm.DB {
	query := db.Instance.
		Limit(pagination.Limit).
		Offset(pagination.Limit).
		Where("date BETWEEN ? AND ?", filters.FromDate, filters.ToDate)

	// Sort
	if sort != nil {
		query.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
	}

	// Categories filter
	if filters.Categories != nil {
		query.Preload("ReceiptItems").Preload("ReceiptItems.Category").Where("category.id IN (?)", filters.Categories)
	}
	// Stores filter
	if len(filters.Stores) > 0 {
		query.Where("store_id IN (?)", filters.Stores)
	}
	// From amount filter
	if filters.FromAmount != nil {
		query.Where("total_purchase_amount > ?", *filters.FromAmount)
	}
	// To amount filter
	if filters.ToAmount != nil {
		query.Where("total_purchase_amount < ?", *filters.ToAmount)
	}

	return query
}

func GetResponseData[T any](models []T) []interface{} {
	var data []interface{}
	for _, model := range models {
		data = append(data, model)
	}

	return data
}

var (
	GorillaLambda *gorillamux.GorillaMuxAdapter
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	// Initialize router
	Router := mux.NewRouter()
	Router.HandleFunc("/receipts", middlewares.SetAuthMiddleware(handler)).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve current User
	user := r.Context().Value(middlewares.CURRENT_USER).(dto.User)

	// Extract query params
	queryParams := r.URL.Query()
	pagination := GetPagination(queryParams)
	filters := GetFilters(queryParams)
	sort := GetSort(queryParams)

	// Build query
	baseQuery := db.Instance.Where("user_id = ?", user.Id)
	query := BuildQuery(baseQuery, pagination, filters, sort)

	// Retrieve receipts
	var dbReceipts []models.Receipt
	queryResult := query.Find(&dbReceipts)
	if queryResult.Error != nil {
		panic(1)
	}

	// Total amount spent
	var total int
	queryResult = query.Select("SUM(total_purchase_amount)").Scan(&total)
	if queryResult.Error != nil {
		panic(1)
	}

	// Convert receipts to interface
	data := GetResponseData[models.Receipt](dbReceipts)

	// Build meta
	paginationMeta := GetPaginationMeta(queryResult, pagination)
	meta := ReceiptsMeta{Pagination: paginationMeta, Total: total}

	// Build response
	response := Response{
		Data: data,
		Meta: meta,
	}

	controllers.JsonResponse(w, &response, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
