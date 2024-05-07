package main

import (
	"context"
	"fmt"
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

type ReceiptDbDto struct {
	ReceiptId    int       `json:"receipt_id"`
	Amount       int       `json:"amount"`
	ReceiptDate  time.Time `json:"receipt_date"`
	StoreName    string    `json:"store_name"`
	CategoryName string    `json:"category_name"`
	Date         time.Time `json:"date"`
}

type Store struct {
	Name string `json:"name"`
}

type Receipt struct {
	ID         int       `json:"id"`
	Amount     string    `json:"amount"`
	Date       time.Time `json:"date"`
	Store      Store     `json:"store"`
	Categories []string  `json:"categories"`
}

type ReceiptsAggregatedByDateListItem struct {
	Date     string    `json:"date"`
	Total    string    `json:"total"`
	Receipts []Receipt `json:"receipts"`
}

type ReceiptsAggregatedByDate = []ReceiptsAggregatedByDateListItem

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
	Router.HandleFunc("/stats/receipts", handler).Methods("GET")
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

	var dbReceipts []ReceiptDbDto
	db.Instance.Model(&models.Receipt{}).
		Select(
			"receipts.id AS receipt_id",
			"receipts.total_purchase_amount AS amount",
			"receipts.date AS receipt_date",
			"stores.name AS store_name",
			"categories.name AS category_name",
			"receipts.date AS receipt_date",
			"DATE(receipts.date) AS date",
		).
		Joins("LEFT JOIN receipt_items ON receipts.id = receipt_items.receipt_id").
		Joins("LEFT JOIN categories ON receipt_items.category_id = categories.id").
		Joins("LEFT JOIN stores ON receipts.store_id = stores.id").
		Where("receipts.date BETWEEN ? AND ?", fromDate, toDate).
		Where("receipts.user_id = ?", user.Id).
		Order("receipts.date DESC").
		Scan(&dbReceipts)

	receiptsMap := make(map[int]Receipt, 0)
	for _, dbReceipt := range dbReceipts {
		receipt, ok := receiptsMap[dbReceipt.ReceiptId]
		if !ok {
			receipt = Receipt{}
			receipt.ID = dbReceipt.ReceiptId
			receipt.Amount = fmt.Sprintf("%.2f", float64(dbReceipt.Amount)/100)
			receipt.Date = dbReceipt.ReceiptDate
			receipt.Store = Store{Name: dbReceipt.StoreName}
			receipt.Categories = make([]string, 0)
		}

		isCategoryUnique := true
		for _, category := range receipt.Categories {
			if category == dbReceipt.CategoryName {
				isCategoryUnique = false
				break
			}
		}

		if isCategoryUnique && dbReceipt.CategoryName != "" {
			receipt.Categories = append(receipt.Categories, dbReceipt.CategoryName)
		}

		receiptsMap[receipt.ID] = receipt
	}

	aggregatedByDate := make(map[string][]Receipt)
	for _, receiptDbDto := range dbReceipts {
		receipt := receiptsMap[receiptDbDto.ReceiptId]
		date := receiptDbDto.Date.Format("2006-01-02")

		isReceiptAdded := false
		for _, receiptsByDate := range aggregatedByDate[date] {
			if receiptsByDate.ID == receipt.ID {
				isReceiptAdded = true
				break
			}
		}

		if !isReceiptAdded {
			aggregatedByDate[date] = append(aggregatedByDate[date], receipt)
		}
	}

	response := make([]ReceiptsAggregatedByDateListItem, 0)
	for date, receipts := range aggregatedByDate {
		total := 0
		for _, receipt := range receipts {
			var dbReceipt ReceiptDbDto
			for _, dbr := range dbReceipts {
				if dbr.ReceiptId == receipt.ID {
					dbReceipt = dbr
					break
				}
			}

			total = total + dbReceipt.Amount
		}

		receiptsAggregatedByDateListItem := ReceiptsAggregatedByDateListItem{}
		receiptsAggregatedByDateListItem.Date = date
		receiptsAggregatedByDateListItem.Total = fmt.Sprintf("%.2f", float64(total)/100)
		receiptsAggregatedByDateListItem.Receipts = receipts
		response = append(response, receiptsAggregatedByDateListItem)
	}

	controllers.JsonResponse(w, response, http.StatusOK)

}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
