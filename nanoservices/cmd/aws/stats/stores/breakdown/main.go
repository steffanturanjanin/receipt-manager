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

type ExpenseByStoreDb struct {
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	Total        int    `json:"total"`
	ReceiptCount int    `json:"receipt_count"`
}

type ExpenseByStore struct {
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	Total        string `json:"total"`
	ReceiptCount int    `json:"receiptCount"`
}

var (
	// Database
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
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
	Router.HandleFunc("/stats/stores/breakdown", handler).Methods("GET")
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

	// Get yearly date range
	fromDate, toDate := GetDateRange()

	var dbExpenses []ExpenseByStoreDb
	dbErr := DB.Model(&models.Store{}).
		Select(
			"stores.tin AS tin",
			"stores.name AS name",
			"IFNULL(SUM(receipts.total_purchase_amount), 0) AS total",
			"COUNT(stores.id) AS receipt_count",
		).
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Where("receipts.user_id = ?", user.Id).
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Group("stores.tin, stores.name").
		Order("total DESC").
		Limit(10).
		Scan(&dbExpenses).
		Error

	if dbErr != nil {
		log.Printf("Error while fetching expenses by category: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	expenses := make([]ExpenseByStore, 0)
	for _, dbExpense := range dbExpenses {
		expense := ExpenseByStore{}
		expense.Tin = dbExpense.Tin
		expense.Name = dbExpense.Name
		expense.Total = fmt.Sprintf("%.2f", float64(dbExpense.Total)/100)
		expense.ReceiptCount = dbExpense.ReceiptCount
		expenses = append(expenses, expense)
	}

	controllers.JsonResponse(w, expenses, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
