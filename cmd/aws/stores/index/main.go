package main

import (
	"log"
	"net/http"
	"os"

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
	authMiddleware := middlewares.SetAuthMiddleware
	queryParamsMiddleware := middlewares.SetQueryParamsMiddleware
	handler := authMiddleware(queryParamsMiddleware(handler))

	// Initialize router
	Router := mux.NewRouter()
	Router.HandleFunc("/stores", middlewares.SetAuthMiddleware(handler)).Methods("GET")
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
	}).Where("receipts.user_id = ?", user.Id)
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
	storesResponse := transport.StoresResponse{}
	storesResponse = storesResponse.FromModels(stores)

	response, err := transport.CreatePaginationResponse(&storesResponse.Items, result.Meta)
	if err != nil {
		log.Printf("Error while building response: %s\n", err.Error())

		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	controllers.JsonResponse(w, &response, http.StatusOK)
}
