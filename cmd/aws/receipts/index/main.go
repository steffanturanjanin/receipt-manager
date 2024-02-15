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

// Response meta object. Includes pagination meta and aggregated
type ReceiptsMeta struct {
	Pagination query.PaginationMeta `json:"pagination"`
	Total      int                  `json:"total"`
}

var (
	GorillaLambda *gorillamux.GorillaMuxAdapter
	DB            *gorm.DB
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}
	DB = db.Instance

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
	queryBuilder := query.NewReceiptQueryBuilder(db.Instance.Where("user_id = ?", user.Id))

	// Extract query params
	sortQuery := middlewares.GetSortQueryParams(r)
	paginationQuery := middlewares.GetPaginationQueryParams(r)
	filterQuery := queryBuilder.GetFilters(r)

	// Execute paginated query
	var receipts []models.Receipt
	queryBuilder = queryBuilder.Filter(filterQuery).Sort(sortQuery)

	result, err := queryBuilder.ExecutePaginatedQuery(&receipts, paginationQuery)
	if err != nil {
		log.Println(err.Error())
		panic(1)
	}

	// Total amount spent
	total, err := queryBuilder.GetTotalPurchaseAmount()
	if err != nil {
		log.Println(err.Error())
		panic(1)
	}

	// Include total in response meta along with pagination meta
	meta := ReceiptsMeta{
		Pagination: result.Meta,
		Total:      *total,
	}

	// Create response object
	response, err := transport.CreatePaginationResponse(result.Data, meta)
	if err != nil {
		log.Printf("Failed to build response: %s\n", err.Error())
		panic(1)
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
