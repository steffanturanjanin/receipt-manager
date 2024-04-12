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

type ExpenseByDate struct {
	Date  string `json:"date"`
	Total string `json:"total"`
}

type ExpenseByDateDB struct {
	Date  string `json:"date"`
	Total int    `json:"total"`
}

var (
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
	Router.HandleFunc("/stats/expenses/breakdown", handler).Methods("GET")
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

	// Execute query
	var expensesDb []ExpenseByDateDB
	dbErr := DB.Model(&models.Receipt{}).
		Select("DATE_FORMAT(receipts.date, '%Y-%m') AS date", "IFNULL(SUM(receipts.total_purchase_amount), 0) AS total").
		Where("receipts.user_id = ?", user.Id).
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Group("DATE_FORMAT(receipts.date, '%Y-%m')").
		Scan(&expensesDb).
		Error

	if dbErr != nil {
		log.Printf("Error while trying to fetch expenses by date: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Initialize response
	expenses := make([]ExpenseByDate, 0)
	for current := fromDate; current.Month() <= toDate.Month() || current.Year() < toDate.Year(); current = current.AddDate(0, 1, 0) {
		expense := ExpenseByDate{
			Date:  current.Format("2006-01"),
			Total: "0.00",
		}
		expenses = append(expenses, expense)
	}

	// Populate response with DB values
	for _, expenseDb := range expensesDb {
		for index, expense := range expenses {
			if expenseDb.Date != expense.Date {
				continue
			}
			expense.Total = fmt.Sprintf("%.2f", float64(expenseDb.Total)/100)
			expenses[index] = expense
			break
		}
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
