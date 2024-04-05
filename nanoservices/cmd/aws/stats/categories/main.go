package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

type CategoryStat struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Total int    `json:"total"`
}

type Category struct {
	ID    int    `json:"id"`
	Color string `json:"color"`
	Name  string `json:"name"`
}

type CategoryStatsResponseItem struct {
	Category Category `json:"category"`
	Total    string   `json:"total"`
}

type CategoryStatsResponse struct {
	Total      string                      `json:"total"`
	Categories []CategoryStatsResponseItem `json:"categories"`
}

var (
	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrMissingFromDateParam = transport.NewBadRequestResponse("Missing fromDate query param")
	ErrMissingToDateParam   = transport.NewBadRequestResponse("Missing toDate query param")
	ErrBadDateFormat        = transport.NewBadRequestResponse("Expected date format is YYYY-MM-dd")
	ErrServiceUnavailable   = transport.NewServiceUnavailableError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	authMiddleware := middlewares.SetAuthMiddleware
	handler := authMiddleware(corsMiddleware(jsonMiddleware(handler)))

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/stats/categories", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	// Retrieve current user
	user := middlewares.GetAuthUser(r)

	// Extract query params
	queryParams := r.URL.Query()

	// Extract date range filter
	fromDateParam := queryParams.Get("fromDate")
	if fromDateParam == "" {
		controllers.JsonResponse(w, ErrMissingFromDateParam, http.StatusBadRequest)
		return
	}
	toDateParam := queryParams.Get("toDate")
	if toDateParam == "" {
		controllers.JsonResponse(w, ErrMissingToDateParam, http.StatusBadRequest)
		return
	}

	// Parse date range params
	timeFormat := "2006-01-02"
	parsedFromDate, err := time.Parse(timeFormat, fromDateParam)
	if err != nil {
		controllers.JsonResponse(w, ErrBadDateFormat, http.StatusBadRequest)
		return
	}
	parsedToDate, err := time.Parse(timeFormat, toDateParam)
	if err != nil {
		controllers.JsonResponse(w, ErrBadDateFormat, http.StatusBadRequest)
		return
	}

	fromDate := time.Date(parsedFromDate.Year(), parsedFromDate.Month(), parsedFromDate.Day(), 0, 0, 0, 0, parsedFromDate.Location())
	toDate := time.Date(parsedToDate.Year(), parsedToDate.Month(), parsedToDate.Day(), 23, 59, 59, 0, parsedToDate.Location())

	var categoriesStats []CategoryStat
	dbErr := db.Instance.Model(&models.Category{}).
		Select("categories.id AS id, categories.name AS name, categories.color AS color, SUM(receipt_items.total_amount) AS total").
		Joins("INNER JOIN receipt_items ON categories.id = receipt_items.category_id").
		Joins("INNER JOIN receipts ON receipt_items.receipt_id = receipts.id").
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Where("receipts.user_id = ?", user.Id).
		Group("categories.id").
		Scan(&categoriesStats).
		Error

	if dbErr != nil {
		log.Printf("Error trying to fetch categories stats: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	var total int
	dbErr = db.Instance.Model(&models.Receipt{}).
		Select("SUM(receipts.total_purchase_amount) AS total").
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Where("receipts.user_id = ?", user.Id).
		Scan(&total).
		Error

	if dbErr != nil {
		log.Printf("Error trying to fetch total spent amount for date range: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	categoryStatsResponse := CategoryStatsResponse{
		Categories: make([]CategoryStatsResponseItem, 0),
		Total:      fmt.Sprintf("%.2f", float64(total)/100),
	}
	for _, categoryStat := range categoriesStats {
		total := fmt.Sprintf("%.2f", float64(categoryStat.Total)/100)
		category := Category{ID: categoryStat.ID, Name: categoryStat.Name, Color: categoryStat.Color}

		categoryStatsResponseItem := CategoryStatsResponseItem{
			Category: category,
			Total:    total,
		}

		categoryStatsResponse.Categories = append(categoryStatsResponse.Categories, categoryStatsResponseItem)
	}

	controllers.JsonResponse(w, &categoryStatsResponse, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
