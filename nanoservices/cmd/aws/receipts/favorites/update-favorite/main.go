package main

import (
	"context"
	"errors"
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
	validation "github.com/steffanturanjanin/receipt-manager/internal/validator"
	"gorm.io/gorm"
)

type FavoriteRequest struct {
	IsFavorite *bool `validate:"required,boolean" json:"isFavorite"`
}

var (
	// Database
	DB *gorm.DB

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Validator
	Validator *validation.Validator

	// Errors
	ErrForbidden          = transport.NewForbiddenError()
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

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/receipts/{id}/favorite", handler).Methods("PATCH")
	GorillaLambda = gorillamux.New(Router)

	// Initialize validator
	Validator = validation.NewDefaultValidator()
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	// Retrieve current user
	user := middlewares.GetAuthUser(r)

	// Receipt id
	receiptId := mux.Vars(r)["id"]

	request := &FavoriteRequest{}
	if err := controllers.ParseBody(request, r); err != nil {
		log.Printf("Error while parsing request: %+v\n", err)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	if err := Validator.GetValidationErrors(request); err != nil {
		controllers.JsonResponse(w, transport.NewValidationError(err), http.StatusUnprocessableEntity)
		return
	}

	// Initialize query
	query := DB.Model(&models.Receipt{}).
		Preload("Store").
		Preload("User").
		Preload("ReceiptItems").
		Preload("ReceiptItems.Category").
		Preload("ReceiptItems.Tax").
		Where("user_id = ?", user.Id)

		// Find Receipt
	var dbReceipt models.Receipt
	if dbErr := query.First(&dbReceipt, receiptId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrForbidden, http.StatusNotFound)
			return
		}

		log.Printf("Error while finding receipt %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	dbReceipt.IsFavorite = *request.IsFavorite
	if dbErr := DB.Save(&dbReceipt).Error; dbErr != nil {
		log.Printf("Error while updating receipt's favorite status: %+v\n", dbErr)
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Build response
	transformer := transport.ReceiptTransformer{}
	receiptResponse := transformer.TransformSingle(dbReceipt)

	// Return response
	controllers.JsonResponse(w, &receiptResponse, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
