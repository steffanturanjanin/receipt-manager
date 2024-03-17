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
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/query"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"gorm.io/gorm"
)

type StoreReceiptsMeta struct {
	Pagination query.PaginationMeta    `json:"pagination"`
	Total      int                     `json:"total"`
	Store      transport.StoreResponse `json:"store"`
}

var (
	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrNotFound           = transport.NewNotFoundError()
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
	Router.HandleFunc("/stores/{id}/receipts", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Store ID path param
	storeIdPathParam := mux.Vars(r)["id"]
	storeId, err := strconv.Atoi(storeIdPathParam)
	if err != nil {
		controllers.JsonResponse(w, ErrNotFound, http.StatusNotFound)
		return
	}

	// Find store
	var store models.Store
	if err := db.Instance.Find(&store, storeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrNotFound, http.StatusNotFound)
			return
		}

		log.Printf("Error while trying to fetch store %d: %s", storeId, err.Error())

		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Auth user
	user := middlewares.GetAuthUser(r)

	// Initialize query builder
	baseQuery := db.Instance.Preload("ReceiptItems").Where("user_id = ?", user.Id).Where("store_id = ?", storeId)
	queryBuilder := query.NewReceiptQueryBuilder(baseQuery)

	// Query params
	sortQuery := middlewares.GetSortQueryParams(r)
	paginationQuery := middlewares.GetPaginationQueryParams(r)
	filterQuery := queryBuilder.GetFilters(r)

	// Execute paginated query
	var receipts []models.Receipt
	paginationResult, err := queryBuilder.Filter(filterQuery).Sort(sortQuery).ExecutePaginatedQuery(&receipts, paginationQuery)
	if err != nil {
		log.Printf("Error while executing paginated query %+v: %s", queryBuilder.Query, err.Error())

		controllers.JsonResponse(w, &ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Total amount spent
	total, err := queryBuilder.GetTotalPurchaseAmount()
	if err != nil {
		log.Printf("Error while executing total purchase amount count %+v: %s", queryBuilder.Query, err.Error())

		controllers.JsonResponse(w, &ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Transform store
	storeTransformer := transport.StoreTransformer{}
	storeResponse := storeTransformer.TransformSingle(store)

	// Build response meta
	meta := StoreReceiptsMeta{
		Pagination: paginationResult.Meta,
		Total:      *total,
		Store:      storeResponse,
	}

	// Transform receipts
	receiptTransformer := transport.BaseReceiptTransformer{}
	data := receiptTransformer.Transform(receipts)

	// Build pagination response
	response, err := transport.CreatePaginationResponse(data, meta)
	if err != nil {
		log.Printf("Error while building pagination response: %s\n", err.Error())

		controllers.JsonResponse(w, &ErrServiceUnavailable, http.StatusServiceUnavailable)
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
