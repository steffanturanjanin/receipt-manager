package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

type FavoriteReceiptStore struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Address  string `json:"address"`
	City     string `json:"city"`
}

type FavoriteReceiptCategory struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type FavoriteReceipt struct {
	ID         int                       `json:"id"`
	Amount     string                    `json:"amount"`
	Date       time.Time                 `json:"date"`
	Store      FavoriteReceiptStore      `json:"store"`
	Categories []FavoriteReceiptCategory `json:"categories"`
}

var (
	// Database
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	//Errors
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
	handler := authMiddleware(corsMiddleware(jsonMiddleware(handler)))

	// Initialize router
	Router := mux.NewRouter()
	Router.HandleFunc("/receipts/favorites", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve current user
	user := middlewares.GetAuthUser(r)

	var dbReceipts []models.Receipt
	dbError := DB.
		Model(&models.Receipt{}).
		Preload("Store", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "location_name", "address", "city")
		}).
		Preload("ReceiptItems.Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "color")
		}).
		Where("user_id = ?", user.Id).
		Where("is_favorite = true").
		Find(&dbReceipts).
		Error

	if dbError != nil {
		log.Printf("Error while trying to fetch favorite receipts: %+v\n", dbError)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	receipts := make([]FavoriteReceipt, 0)

	for _, dbReceipt := range dbReceipts {
		categories := make([]FavoriteReceiptCategory, 0)

		for _, dbReceiptItem := range dbReceipt.ReceiptItems {
			if dbReceiptItem.Category == nil {
				continue
			}

			isCategoryPresent := false

			for _, category := range categories {
				if category.ID == int(dbReceiptItem.Category.ID) {
					isCategoryPresent = true
					break
				}
			}

			if !isCategoryPresent {
				category := FavoriteReceiptCategory{}
				category.ID = int(dbReceiptItem.Category.ID)
				category.Name = dbReceiptItem.Category.Name
				category.Color = dbReceiptItem.Category.Color
				categories = append(categories, category)
			}
		}

		receipt := FavoriteReceipt{
			ID:         int(dbReceipt.ID),
			Date:       *dbReceipt.Date,
			Amount:     fmt.Sprintf("%.2f", float64(*dbReceipt.TotalPurchaseAmount)/100),
			Categories: categories,
			Store: FavoriteReceiptStore{
				ID:       int(dbReceipt.Store.ID),
				Name:     dbReceipt.Store.Name,
				Location: dbReceipt.Store.LocationName,
				Address:  dbReceipt.Store.Address,
				City:     dbReceipt.Store.City,
			},
		}

		receipts = append(receipts, receipt)
	}

	controllers.JsonResponse(w, receipts, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
