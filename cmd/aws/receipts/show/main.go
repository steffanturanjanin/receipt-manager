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

	//Errors
	ErrReceiptNotFound    = transport.NewNotFoundError()
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
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
	// Get Auth user
	user := middlewares.GetAuthUser(r)

	// Extract id from path params
	pathParams := mux.Vars(r)
	receiptId, _ := strconv.ParseInt(pathParams["id"], 10, 64)

	// Initialize query
	query := db.Instance.Where("user_id = ?", user.Id)

	// Find Receipt
	var dbReceipt models.Receipt
	if dbErr := query.First(&dbReceipt, receiptId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrReceiptNotFound, http.StatusNotFound)
			return
		}

		log.Printf("Error while finding receipt %d: %s", receiptId, dbErr.Error())
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
	}

	// Build response
	receiptResponse := transport.ReceiptResponse{}
	receiptResponse = receiptResponse.FromModel(dbReceipt)

	// Return response
	controllers.JsonResponse(w, &receiptResponse, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
