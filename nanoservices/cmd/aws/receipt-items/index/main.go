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

type ReceiptItemDb struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Store    string    `json:"store"`
	Amount   int       `json:"amount"`
	Unit     string    `json:"unit"`
	Quantity int       `json:"quantity"`
}

type ReceiptItem struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Store    string    `json:"store"`
	Amount   string    `json:"amount"`
	Unit     string    `json:"unit"`
	Quantity int       `json:"quantity"`
}

var (
	// Database
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrMissingSearchText  = transport.NewBadRequestResponse("Missing required 'searchText' query parameter")
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
	Router.HandleFunc("/receipt-items", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	// Retrieve current user
	user := middlewares.GetAuthUser(r)

	queryParams := r.URL.Query()
	searchText := queryParams.Get("searchText")

	if searchText == "" {
		controllers.JsonResponse(w, ErrMissingSearchText, http.StatusBadRequest)
		return
	}

	var dbReceiptItems []ReceiptItemDb
	dbErr := DB.Model(&models.ReceiptItem{}).
		Select(
			"receipt_items.id AS id",
			"receipt_items.name AS name",
			"receipt_items.single_amount AS amount",
			"receipt_items.unit AS unit",
			"receipt_items.quantity AS quantity",
			"receipts.date AS date",
			"stores.name AS store",
		).
		Joins("INNER JOIN receipts ON receipt_items.receipt_id = receipts.id").
		Joins("INNER JOIN stores ON receipts.store_id = stores.id").
		Where("receipts.user_id = ?", user.Id).
		Where("receipt_items.name LIKE ? OR stores.name LIKE ?", "%"+searchText+"%", "%"+searchText+"%").
		Order("receipt_items.name ASC, receipts.date DESC").
		Scan(&dbReceiptItems).
		Error

	if dbErr != nil {
		log.Printf("Error while fetching receipt items: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	receiptItems := make([]ReceiptItem, 0)
	for _, dbReceiptItem := range dbReceiptItems {
		receiptItem := ReceiptItem{}
		receiptItem.ID = dbReceiptItem.ID
		receiptItem.Name = dbReceiptItem.Name
		receiptItem.Date = dbReceiptItem.Date
		receiptItem.Store = dbReceiptItem.Store
		receiptItem.Unit = dbReceiptItem.Unit
		receiptItem.Quantity = dbReceiptItem.Quantity
		receiptItem.Amount = fmt.Sprintf("%.2f", float64(dbReceiptItem.Amount)/100)
		receiptItems = append(receiptItems, receiptItem)
	}

	controllers.JsonResponse(w, receiptItems, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
