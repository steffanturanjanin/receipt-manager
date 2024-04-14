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
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

type Store struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	City     string `json:"city"`
	Address  string `json:"address"`
}

var (
	// Database
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrMissingSearchText  = transport.NewBadRequestResponse("Missing required 'searchText' query parameter")
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	} else {
		DB = db.Instance
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

	queryParams := r.URL.Query()
	searchText := queryParams.Get("searchText")

	if searchText == "" {
		controllers.JsonResponse(w, ErrMissingSearchText, http.StatusBadRequest)
		return
	}

	var dbStores []models.Store
	dbErr := DB.Select("stores.id AS id", "name", "location_name", "city", "address").
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Where("receipts.user_id = ?", user.Id).
		Where("name LIKE ? OR location_name LIKE ? OR city LIKE ? OR address LIKE ?", "%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%").
		Find(&dbStores).
		Error

	if dbErr != nil {
		log.Printf("Error while fetching stores: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	stores := make([]Store, 0)
	for _, dbStore := range dbStores {
		store := Store{}
		store.ID = int(dbStore.ID)
		store.Name = dbStore.Name
		store.Location = dbStore.LocationName
		store.Address = dbStore.Address
		store.City = dbStore.City
		stores = append(stores, store)
	}

	controllers.JsonResponse(w, stores, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
