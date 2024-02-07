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
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

var (
	// Router
	gorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrReceiptNotFound = transport.ErrorResponse{
		Error: "receipt not found",
		Code:  404,
	}
	ErrReceiptForbidden = transport.ErrorResponse{
		Error: "forbidden",
		Code:  403,
	}
)

func init() {
	// Initialize database
	err := db.InitializeDB()
	if err != nil {
		os.Exit(1)
	}

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/receipts/{id}", middlewares.SetAuthMiddleware(handler)).Methods("DELETE")
	gorillaLambda = gorillamux.New(router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middlewares.CURRENT_USER).(dto.User)
	receiptId, _ := strconv.Atoi(mux.Vars(r)["id"])

	var dbReceipt models.Receipt
	// Check if receipt exists
	if dbErr := db.Instance.Where(&dbReceipt).First(receiptId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			// Not Found 404
			controllers.JsonResponse(w, ErrReceiptNotFound, http.StatusNotFound)
			return
		}
		// Unknown error. Lambda cannot process this error.
		panic(1)
	}

	// Check if user has permission to delete receipt
	if dbReceipt.UserID != &user.Id {
		// Forbidden 403
		controllers.JsonResponse(w, ErrReceiptForbidden, http.StatusForbidden)
		return
	}

	// Delete Receipt
	if dbError := db.Instance.Delete(&dbReceipt); dbError != nil {
		// Unknown error. Lambda cannot process this error.
		panic(1)
	}

	// Set No Content 204
	w.WriteHeader(http.StatusNoContent)
}

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(lambdaHandler)
}
