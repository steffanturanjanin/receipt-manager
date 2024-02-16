package main

import (
	"context"
	"errors"
	"log"
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
)

var (
	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrNotFound           = transport.NewNotFoundError()
	ErrForbidden          = transport.NewForbiddenError()
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/receipts/{id}", middlewares.SetAuthMiddleware(handler)).Methods("DELETE")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetAuthUser(r)
	receiptId, _ := strconv.Atoi(mux.Vars(r)["id"])

	var dbReceipt models.Receipt
	// Check if receipt exists
	if dbErr := db.Instance.First(&dbReceipt, receiptId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			// Not Found 404
			controllers.JsonResponse(w, ErrNotFound, http.StatusNotFound)
			return
		}
		// Unknown error. Lambda cannot process this error.
		log.Printf("Error while trying to fetch receipt %d: %s\n", receiptId, dbErr.Error())
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
	}

	// Check if user has permission to delete receipt
	if *dbReceipt.UserID != user.Id {
		// Forbidden 403
		controllers.JsonResponse(w, ErrForbidden, http.StatusForbidden)
		return
	}

	// Delete Receipt
	if dbErr := db.Instance.Delete(&dbReceipt).Error; dbErr != nil {
		// Unknown error. Lambda cannot process this error.
		log.Printf("Error while trying to delete receipt %d: %s\n", receiptId, dbErr.Error())
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
	}

	// Set No Content 204
	w.WriteHeader(http.StatusNoContent)
}

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(lambdaHandler)
}
