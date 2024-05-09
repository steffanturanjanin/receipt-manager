package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
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
	"gorm.io/gorm"
)

type MostPopularReceiptItemDB struct {
	Name         string `json:"name"`
	ReceiptCount int    `json:"receipt_count"`
	Total        int    `json:"total"`
}

type MostPopularReceiptItem struct {
	Name         string `json:"name"`
	ReceiptCount int    `json:"receiptCount"`
	Total        string `json:"total"`
}

type MostPopularStoreDB struct {
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	Address      string `json:"address"`
	City         string `json:"city"`
	ReceiptCount int    `json:"receipt_count"`
	Total        int    `json:"total"`
}

type MostPopularStore struct {
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	Address      string `json:"address"`
	City         string `json:"city"`
	ReceiptCount int    `json:"receiptCount"`
	Total        string `json:"total"`
	Percent      int    `json:"percent"`
}

type CategoryStats struct {
	ID                      int                      `json:"id"`
	Name                    string                   `json:"name"`
	Total                   string                   `json:"total"`
	MostPopularReceiptItems []MostPopularReceiptItem `json:"mostPopularReceiptItems"`
	MostPopularStores       []MostPopularStore       `json:"mostPopularStores"`
}

var (
	// Database
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrNotFound           = transport.NewNotFoundError()
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	} else {
		DB = db.Instance
	}

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	authMiddleware := middlewares.SetAuthMiddleware
	handler := authMiddleware(corsMiddleware(jsonMiddleware(handler)))

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/stats/categories/{id}", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func GetDateRange() (time.Time, time.Time) {
	currentDate := time.Now()
	firstDateOfCurrentMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	startDate := firstDateOfCurrentMonth.AddDate(0, -11, 0)

	fromDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	toDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 23, 59, 59, 0, currentDate.Location())

	return fromDate, toDate
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	// Retrieve current user
	user := middlewares.GetAuthUser(r)
	// Get category id from url path segment
	categoryId := mux.Vars(r)["id"]

	// Get yearly date range
	fromDate, toDate := GetDateRange()

	var category models.Category
	if dbError := DB.Select("id", "name").Where("id = ?", categoryId).First(&category).Error; dbError != nil {
		if errors.Is(dbError, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrNotFound, http.StatusNotFound)
			return
		}

		log.Printf("Error while trying to fetch category: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	var total int
	dbError := DB.
		Model(&models.Category{}).
		Select("IFNULL(SUM(receipts.total_purchase_amount), 0) AS total").
		Joins("INNER JOIN receipt_items ON categories.id = receipt_items.category_id").
		Joins("INNER JOIN receipts ON receipt_items.receipt_id = receipts.id").
		Where("receipts.user_id = ?", user.Id).
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Where("categories.id = ?", categoryId).
		Group("categories.id").
		Scan(&total).
		Error

	if dbError != nil {
		log.Printf("Error while trying to fetch totals for category: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	var dbMostPopularReceiptItems []MostPopularReceiptItemDB
	dbError = DB.
		Model(&models.ReceiptItem{}).
		Select(
			"receipt_items.name AS name",
			"COUNT(receipt_items.name) AS receipt_count",
			"IFNULL(SUM(receipt_items.total_amount), 0) AS total",
		).
		Joins("INNER JOIN receipts ON receipt_items.receipt_id = receipts.id").
		Where("receipts.user_id = ?", user.Id).
		Where("receipt_items.category_id = ?", categoryId).
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Group("receipt_items.name").
		Order("total DESC").
		Scan(&dbMostPopularReceiptItems).
		Error

	if dbError != nil {
		log.Printf("Error while fetching most popular receipt items for category: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	var dbMostPopularStores []MostPopularStoreDB
	dbError = DB.
		Model(&models.Store{}).
		Select(
			"stores.tin AS tin",
			"stores.name AS name",
			"stores.location_name AS location",
			"stores.address AS address",
			"stores.city AS city",
			"COUNT(stores.id) AS receipt_count",
			"IFNULL(SUM(receipts.total_purchase_amount), 0) AS total",
		).
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Joins("INNER JOIN receipt_items ON receipts.id = receipt_items.receipt_id").
		Where("receipts.user_id = ?", user.Id).
		Where("receipt_items.category_id = ?", categoryId).
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Group("stores.id").
		Order("total DESC, receipt_count DESC, stores.name DESC").
		Scan(&dbMostPopularStores).
		Error

	if dbError != nil {
		log.Printf("Error while fetching most popular stores for category: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	mostPopularReceiptItems := make([]MostPopularReceiptItem, 0)
	for _, dbMostPopularReceiptItem := range dbMostPopularReceiptItems {
		receiptItem := MostPopularReceiptItem{}
		receiptItem.Name = dbMostPopularReceiptItem.Name
		receiptItem.ReceiptCount = dbMostPopularReceiptItem.ReceiptCount
		receiptItem.Total = fmt.Sprintf("%.2f", float64(dbMostPopularReceiptItem.Total)/100)
		mostPopularReceiptItems = append(mostPopularReceiptItems, receiptItem)
	}

	mostPopularStores := make([]MostPopularStore, 0)
	for _, dbMostPopularStore := range dbMostPopularStores {
		store := MostPopularStore{}
		store.Tin = dbMostPopularStore.Tin
		store.Name = dbMostPopularStore.Name
		store.Location = dbMostPopularStore.Location
		store.Address = dbMostPopularStore.Address
		store.City = dbMostPopularStore.City
		store.ReceiptCount = dbMostPopularStore.ReceiptCount
		store.Total = fmt.Sprintf("%.2f", float64(dbMostPopularStore.Total)/100)
		store.Percent = int(math.Round(float64(dbMostPopularStore.Total) / float64(total) * 100))
		mostPopularStores = append(mostPopularStores, store)
	}

	categoryStats := CategoryStats{}
	categoryStats.ID = int(category.ID)
	categoryStats.Name = category.Name
	categoryStats.Total = fmt.Sprintf("%.2f", float64(total)/100)
	categoryStats.MostPopularReceiptItems = mostPopularReceiptItems
	categoryStats.MostPopularStores = mostPopularStores

	controllers.JsonResponse(w, categoryStats, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
