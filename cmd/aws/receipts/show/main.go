package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

var (
	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	//Errors
	ErrReceiptNotFound = transport.NewNotFoundError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	// Initialize router
	Router := mux.NewRouter()
	Router.HandleFunc("/receipts/{id}", middlewares.SetAuthMiddleware(handler)).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract id from path params
	pathParams := mux.Vars(r)
	receiptId, _ := strconv.ParseInt(pathParams["id"], 10, 64)

	// Find Receipt
	var dbReceipt models.Receipt
	if dbErr := db.Instance.Find(&dbReceipt, receiptId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrReceiptNotFound, http.StatusNotFound)
			return
		}

		panic(1)
	}

	// Return response
	controllers.JsonResponse(w, &dbReceipt, http.StatusOK)
}