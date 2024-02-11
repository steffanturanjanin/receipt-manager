package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	validation "github.com/steffanturanjanin/receipt-manager/internal/validator"
)

type UpdateReceiptItemCategoryRequest struct {
	CategoryId uint `validate:"required" json:"categoryId"`
}

var (
	// Router
	gorillaLambda *gorillamux.GorillaMuxAdapter

	// Validator
	validator *validation.Validator

	//Errors
	ErrReceiptItemNotFound = transport.NewNotFoundError()
	ErrCategoryNotFound    = transport.NewNotFoundError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/receipt-items/{id}", middlewares.SetAuthMiddleware(handler)).Methods("PATCH")
	gorillaLambda = gorillamux.New(router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to struct
	updateReceiptItemCategoryRequest := &UpdateReceiptItemCategoryRequest{}
	if err := controllers.ParseBody(updateReceiptItemCategoryRequest, r); err != nil {
		panic(1)
	}

	// Validate request
	// If failed return 422 Unprocessed Entity with error map
	if err := validator.GetValidationErrors(updateReceiptItemCategoryRequest); err != nil {
		controllers.JsonResponse(w, transport.NewValidationError(err), http.StatusUnprocessableEntity)
		return
	}

	pathParams := mux.Vars(r)
	receiptItemId, err := strconv.ParseInt(pathParams["id"], 10, 64)
	if err != nil {
		panic(1)
	}

	// Hydrate Receipt Item
	// If Not Found - return 404 HTTP status code
	var dbReceiptItem models.ReceiptItem
	if dbErr := db.Instance.Find(&dbReceiptItem, receiptItemId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrReceiptItemNotFound, http.StatusNotFound)
			return
		}

		panic(1)
	}

	// Hydrate Category
	// If Not Found - return 404 HTTP status code
	var dbCategory models.Category
	if dbErr := db.Instance.Find(&dbCategory, updateReceiptItemCategoryRequest.CategoryId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrCategoryNotFound, http.StatusNotFound)
			return
		}

		panic(1)
	}

	// Update Receipt Item Category
	*dbReceiptItem.CategoryID = updateReceiptItemCategoryRequest.CategoryId
	if db.Instance.Save(&dbReceiptItem).Error != nil {
		panic(1)
	}

	controllers.JsonResponse(w, &dbReceiptItem, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
