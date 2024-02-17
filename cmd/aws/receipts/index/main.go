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
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/query"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

// Response meta object. Includes pagination meta and aggregated
type ReceiptsMeta struct {
	Pagination query.PaginationMeta `json:"pagination"`
	Total      int                  `json:"total"`
}

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
	Router.HandleFunc("/receipts", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve current User
	user := middlewares.GetAuthUser(r)

	// Initialize query builder
	baseQuery := db.Instance.Preload("Store").Preload("ReceiptItems").Where("user_id = ?", user.Id)
	queryBuilder := query.NewReceiptQueryBuilder(baseQuery)

	// Extract query params
	sortQuery := middlewares.GetSortQueryParams(r)
	paginationQuery := middlewares.GetPaginationQueryParams(r)
	filterQuery := queryBuilder.GetFilters(r)

	// Execute paginated query
	var receipts []models.Receipt
	result, err := queryBuilder.Filter(filterQuery).Sort(sortQuery).ExecutePaginatedQuery(&receipts, paginationQuery)
	if err != nil {
		log.Printf("Error while executing paginated query %+v: %s", queryBuilder.Query, err.Error())

		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Transformed receipts response
	receiptsResponse := transport.ReceiptsResponse{}
	receiptsResponse = receiptsResponse.FromModels(receipts)

	// Total amount spent
	total, err := queryBuilder.GetTotalPurchaseAmount()
	if err != nil {
		log.Printf("Error while executing total purchase amount count %+v: %s", queryBuilder.Query, err.Error())

		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Include total in response meta along with pagination meta
	meta := ReceiptsMeta{
		Pagination: result.Meta,
		Total:      *total,
	}

	// Create response object
	response, err := transport.CreatePaginationResponse(&receiptsResponse.Items, meta)
	if err != nil {
		log.Printf("Error while building response object: %s", err.Error())

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
