package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"gorm.io/gorm"
)

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
	if err := database.InitializeDB(); err != nil {
		os.Exit(1)
	} else {
		DB = database.Instance
	}

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	authMiddleware := middlewares.SetAuthMiddleware
	handler := authMiddleware(corsMiddleware(jsonMiddleware(handler)))

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/categories", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	if dbErr := DB.Find(&categories).Error; dbErr != nil {
		log.Printf("Error while fetching categories %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	transformer := transport.CategoryTransformer{}
	response := transformer.Transform(categories)

	controllers.JsonResponse(w, response, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
