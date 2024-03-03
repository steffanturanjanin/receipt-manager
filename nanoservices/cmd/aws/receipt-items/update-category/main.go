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
	ut "github.com/go-playground/universal-translator"
	v "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"github.com/steffanturanjanin/receipt-manager/internal/validator"
)

type UpdateReceiptItemCategoryRequest struct {
	CategoryId uint `validate:"required,category_exists" json:"categoryId"`
}

var (
	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Validator
	V *validator.Validator

	//Errors
	ErrReceiptItemNotFound = transport.NewNotFoundError()
	ErrCategoryNotFound    = transport.NewNotFoundError()
	ErrServiceUnavailable  = transport.NewServiceUnavailableError()
)

func validateCategoryExistence(fl v.FieldLevel) bool {
	categoryId := fl.Field().Uint()
	if dbErr := db.Instance.Find(&models.Category{}, categoryId).Error; dbErr != nil {
		return false
	}

	return true
}

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	// Initialize router
	Router := mux.NewRouter()
	Router.HandleFunc("/receipt-items/{id}", middlewares.SetAuthMiddleware(handler)).Methods("PATCH")
	GorillaLambda = gorillamux.New(Router)

	// Initialize validator
	V = validator.NewDefaultValidator()
	V.Validator.RegisterValidation("category_exists", validateCategoryExistence)
	V.Validator.RegisterTranslation("category_exists", V.Translator, func(ut ut.Translator) error {
		return ut.Add("category_exists", "Category with id {0} does not exist.", true)
	}, func(ut ut.Translator, fe v.FieldError) string {
		value, _ := fe.Value().(string)
		t, _ := ut.T("category_exists", value)
		return t
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to struct
	updateReceiptItemCategoryRequest := &UpdateReceiptItemCategoryRequest{}
	if err := controllers.ParseBody(updateReceiptItemCategoryRequest, r); err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
	}

	// Validate request
	// If failed return 422 Unprocessed Entity with error map
	if err := V.GetValidationErrors(updateReceiptItemCategoryRequest); err != nil {
		controllers.JsonResponse(w, transport.NewValidationError(err), http.StatusUnprocessableEntity)
		return
	}

	// Extract receipt item {id} path parameter
	pathParams := mux.Vars(r)
	receiptItemId, err := strconv.ParseInt(pathParams["id"], 10, 64)
	if err != nil {
		log.Printf("Error converting path param: %s\n", err.Error())
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
	}

	// Hydrate Receipt Item
	// If Not Found - return 404 HTTP status code
	var dbReceiptItem models.ReceiptItem
	if dbErr := db.Instance.Find(&dbReceiptItem, receiptItemId).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			controllers.JsonResponse(w, ErrReceiptItemNotFound, http.StatusNotFound)
			return
		}

		log.Printf("Error fetching receipt item %+d: %s\n", receiptItemId, err.Error())
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
	}

	// Update Receipt Item Category
	*dbReceiptItem.CategoryID = updateReceiptItemCategoryRequest.CategoryId
	if db.Instance.Save(&dbReceiptItem).Error != nil {
		log.Printf("Error updating receipt item %+d: %s\n", receiptItemId, err.Error())
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
	}

	// Transform response
	transformer := transport.ReceiptItemTransformer{}
	response := transformer.TransformSingle(dbReceiptItem)

	controllers.JsonResponse(w, &response, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
