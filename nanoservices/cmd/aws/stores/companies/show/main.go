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
	"gorm.io/gorm"
)

type LocationDb struct {
	LocationID   string `json:"location_id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
	City         string `json:"city"`
	ReceiptCount int    `json:"receipt_count"`
	Amount       int    `json:"amount"`
}

type Location struct {
	LocationID   string `json:"locationId"`
	LocationName string `json:"locationName"`
	Address      string `json:"address"`
	City         string `json:"city"`
	ReceiptCount int    `json:"receiptCount"`
	Amount       string `json:"amount"`
}

type LocationTotals struct {
	ReceiptCount int `json:"receipt_count"`
	Amount       int `json:"amount"`
}

type Locations struct {
	Data         []Location `json:"data"`
	Total        string     `json:"total"`
	ReceiptCount int        `json:"receiptCount"`
}

type LocationExpenseDb struct {
	ID           int       `json:"id"`
	LocationName string    `json:"location_name"`
	Date         time.Time `json:"date"`
	Amount       int       `json:"amount"`
}

type LocationExpense struct {
	ID           int       `json:"id"`
	LocationName string    `json:"locationName"`
	Date         time.Time `json:"date"`
	Amount       string    `json:"amount"`
}

type Company struct {
	Tin       string            `json:"tin"`
	Name      string            `json:"name"`
	Locations Locations         `json:"locations"`
	Expenses  []LocationExpense `json:"expenses"`
}

var (
	// Database
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
	ErrNotFound           = transport.NewNotFoundError()
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
	Router.HandleFunc("/stores/companies/{tin}", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	// Retrieve current user
	user := middlewares.GetAuthUser(r)

	// Extract tin from route path
	tin := mux.Vars(r)["tin"]

	// Store
	var storeDb *models.Store
	dbError := DB.
		Select("stores.tin AS tin", "stores.name AS name").
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Where("tin = ?", tin).
		Where("receipts.user_id = ?", user.Id).
		First(&storeDb).
		Error

	if dbError != nil {
		log.Printf("Error while fetching store: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	if storeDb == nil {
		controllers.JsonResponse(w, ErrNotFound, http.StatusNotFound)
		return
	}

	// Location data
	var locationsDb []LocationDb
	dbError = DB.
		Model(&models.Store{}).
		Select(
			"stores.location_id AS location_id",
			"stores.location_name AS location_name",
			"stores.address AS address",
			"stores.city AS city",
			"COUNT(stores.location_id) AS receipt_count",
			"IFNULL(SUM(receipts.total_purchase_amount), 0) AS amount",
		).
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Where("stores.tin = ?", tin).
		Where("receipts.user_id = ?", user.Id).
		Group("stores.location_id, stores.id").
		Order("IFNULL(SUM(receipts.total_purchase_amount), 0) DESC, stores.location_name ASC").
		Scan(&locationsDb).
		Error

	if dbError != nil {
		log.Printf("Error while trying to fetch locations data: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Location totals
	var locationTotals LocationTotals
	dbError = DB.
		Model(&models.Store{}).
		Select(
			"COUNT(stores.location_id) AS receipt_count",
			"SUM(receipts.total_purchase_amount) AS amount",
		).
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Where("stores.tin = ?", tin).
		Where("receipts.user_id = ?", user.Id).
		Scan(&locationTotals).
		Error

	if dbError != nil {
		log.Printf("Error while trying to fetch locations totals: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Expenses by locations
	var expensesDb []LocationExpenseDb
	dbError = DB.
		Model(&models.Store{}).
		Select(
			"receipts.id AS id",
			"receipts.date AS date",
			"receipts.total_purchase_amount AS amount",
			"stores.location_name AS location_name",
		).
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Where("stores.tin = ?", tin).
		Where("receipts.user_id = ?", user.Id).
		Order("receipts.date DESC, stores.location_name ASC").
		Scan(&expensesDb).
		Error

	if dbError != nil {
		log.Printf("Error while fetching store expenses: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	store := Company{
		Tin:  storeDb.Tin,
		Name: storeDb.Name,
		Locations: Locations{
			Data:         make([]Location, 0),
			ReceiptCount: locationTotals.ReceiptCount,
			Total:        fmt.Sprintf("%.2f", float64(locationTotals.Amount)/100),
		},
		Expenses: make([]LocationExpense, 0),
	}

	for _, locationDb := range locationsDb {
		location := Location{}
		location.LocationID = locationDb.LocationID
		location.LocationName = locationDb.LocationName
		location.Address = locationDb.Address
		location.City = locationDb.City
		location.ReceiptCount = locationDb.ReceiptCount
		location.Amount = fmt.Sprintf("%.2f", float64(locationDb.Amount)/100)
		store.Locations.Data = append(store.Locations.Data, location)
	}

	for _, expenseDb := range expensesDb {
		expense := LocationExpense{}
		expense.ID = expenseDb.ID
		expense.LocationName = expenseDb.LocationName
		expense.Date = expenseDb.Date
		expense.Amount = fmt.Sprintf("%.2f", float64(expenseDb.Amount)/100)
		store.Expenses = append(store.Expenses, expense)
	}

	controllers.JsonResponse(w, store, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
