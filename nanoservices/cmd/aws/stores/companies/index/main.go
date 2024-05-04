package main

import (
	"context"
	"fmt"
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
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"gorm.io/gorm"
)

type CompanyDb struct {
	ID           int    `json:"id"`
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	Total        int    `json:"total"`
	ReceiptCount int    `json:"receipt_count"`
}

type Company struct {
	ID           int    `json:"id"`
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	Total        string `json:"total"`
	ReceiptCount int    `json:"receiptCount"`
}

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
	Router.HandleFunc("/stores/companies", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Auth user
	user := middlewares.GetAuthUser(r)

	queryParams := r.URL.Query()
	searchText := queryParams.Get("searchText")

	var dbCompanies []CompanyDb
	query := DB.
		Model(&models.Store{}).
		Select(
			"stores.id AS id",
			"stores.tin AS tin",
			"stores.name AS name",
			"COUNT(stores.id) AS receipt_count",
			"SUM(receipts.total_purchase_amount) AS total",
		).
		Joins("INNER JOIN receipts ON stores.id = receipts.store_id").
		Where("receipts.user_id = ?", user.Id).
		Group("stores.tin, stores.id").
		Order("stores.name ASC")

	if searchText != "" {
		query.Where("stores.name LIKE ?", "%"+searchText+"%")
	}

	dbErr := query.Scan(&dbCompanies).Error

	if dbErr != nil {
		log.Printf("Error while fetching companies: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	companies := make([]Company, 0)
	for _, dbCompany := range dbCompanies {
		company := Company{}
		company.ID = dbCompany.ID
		company.Name = dbCompany.Name
		company.Tin = dbCompany.Tin
		company.ReceiptCount = dbCompany.ReceiptCount
		company.Total = fmt.Sprintf("%.2f", float64(dbCompany.Total)/100)
		companies = append(companies, company)
	}

	controllers.JsonResponse(w, companies, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
