package main

import (
	"context"
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
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"gorm.io/gorm"

	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
)

type ProfileDb struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	ReceiptCount int       `json:"receipt_count"`
	RegisteredAt time.Time `json:"registered_at"`
}

type Profile struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	ReceiptCount int       `json:"receiptCount"`
	RegisteredAt time.Time `json:"registeredAt"`
}

var (
	// Database instance
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		log.Printf("Error while trying to initialize the database: %+v\n", err)
		os.Exit(1)
	} else {
		DB = db.Instance
	}

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	authMiddleware := middlewares.SetAuthMiddleware
	handler := authMiddleware(corsMiddleware(jsonMiddleware(handler)))

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/auth/me", handler).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	// Get Auth user
	user := middlewares.GetAuthUser(r)

	var profileDb ProfileDb
	dbErr := DB.Model(&models.User{}).
		Select(
			"users.id AS id",
			"users.first_name AS first_name",
			"users.last_name AS last_name",
			"users.email AS email",
			"users.created_at AS registered_at",
			"COUNT(receipts.id) AS receipt_count",
		).
		Joins("INNER JOIN receipts ON users.id = receipts.user_id").
		Where("users.id = ?", user.Id).
		Group("users.id").
		Scan(&profileDb).
		Error

	if dbErr != nil {
		log.Printf("Error while trying to count user's receipts: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	profile := Profile{}
	profile.ID = profileDb.ID
	profile.FirstName = profileDb.FirstName
	profile.LastName = profileDb.LastName
	profile.Email = profileDb.Email
	profile.RegisteredAt = profileDb.RegisteredAt
	profile.ReceiptCount = profileDb.ReceiptCount

	// Return profile response
	controllers.JsonResponse(w, profile, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
