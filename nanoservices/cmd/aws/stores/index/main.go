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
	"gorm.io/gorm"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/query"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

var (
	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
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
	queryParamsMiddleware := middlewares.SetQueryParamsMiddleware
	handler := authMiddleware(corsMiddleware(queryParamsMiddleware(jsonMiddleware(handler))))

	// Initialize router
	Router := mux.NewRouter()
	Router.HandleFunc("/stores", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Auth user
	user := middlewares.GetAuthUser(r)

	// Query params
	paginationQuery := middlewares.GetPaginationQueryParams(r)

	// Query builder
	baseQuery := db.Instance.Preload("Receipts", func(query *gorm.DB) *gorm.DB {
		return query.Select("user_id")
	}).Where("receipts.user_id = ?", user.Id).Order("name desc")

	queryBuilder := query.NewStoreQueryBuilder(baseQuery)

	// Execute paginated query
	var stores []models.Store
	result, err := queryBuilder.ExecutePaginatedQuery(&stores, paginationQuery)
	if err != nil {
		log.Printf("Error trying to execute paginated query: %s\n", err.Error())

		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Build response
	transformer := transport.StoreTransformer{}
	storesResponse := transformer.Transform(stores)

	response, err := transport.CreatePaginationResponse(&storesResponse, result.Meta)
	if err != nil {
		log.Printf("Error while building response: %s\n", err.Error())

		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	controllers.JsonResponse(w, &response, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
