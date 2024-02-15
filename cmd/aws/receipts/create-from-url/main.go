package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	validation "github.com/steffanturanjanin/receipt-manager/internal/validator"
)

/**
1) Check if user is logged-in
2) Validate submitted URL:
- check if URL is present
- check if URL starts with "https://suf.purs.gov.rs"
- check if URL has ?vl= query string

3) Extract ?vl= query string parameter from URL
4) Check if database has an entry with this user_id and vl query string
- if there is an entry, return validation error: Receipt is already submitted.

5) Add new entry to the receipts table:
- user_id
- vl query string parameter
- status: pending

6) Write URL to SQS queue for further processing
7) Return 201 Created HTTP response code
*/

type ReceiptUrlRequest struct {
	Url string `validate:"required,url,url_host=suf.purs.gov.rs,url_query_params=vl" json:"url"`
}

const PENDING_RECEIPTS_QUEUE = "pending_receipts"

var (
	// Execution
	env = os.Getenv("ENVIRONMENT")

	// AWS SQS
	session               *aws_session.Session
	client                *sqs.SQS
	pendingReceiptsSqsUrl *string

	// Router
	gorillaLambda *gorillamux.GorillaMuxAdapter

	// Validator
	validator *validation.Validator

	// Errors
	ErrReceiptAlreadyScanned = transport.NewBadRequestResponse(errors.New("receipt already scanned"))
	ErrCannotCreateReceipt   = transport.NewBadRequestResponse(errors.New("receipt cannot be created"))
)

func init() {
	// Initialize database
	err := database.InitializeDB()
	if err != nil {
		os.Exit(1)
	}

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/receipts", middlewares.SetAuthMiddleware(handler)).Methods("POST")
	gorillaLambda = gorillamux.New(router)

	// Initialize AWS session and SQS client
	sessionOptions := aws_session.Options{
		Config:            aws.Config{Endpoint: aws.String("http://docker.for.mac.localhost:9324")}, // TODO: configure this to run only in DEV environment
		SharedConfigState: aws_session.SharedConfigDisable,
	}
	session = aws_session.Must(aws_session.NewSessionWithOptions(sessionOptions))
	client = sqs.New(session)

	// Initialize SQS urls
	if urlResult, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(PENDING_RECEIPTS_QUEUE),
	}); err != nil {
		panic(1)
	} else {
		pendingReceiptsSqsUrl = urlResult.QueueUrl
	}

	// Initialize validator
	validator = validation.NewDefaultValidator()
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middlewares.CURRENT_USER).(dto.User)

	// Parse request body to struct
	receiptUrlRequest := &ReceiptUrlRequest{}
	if err := controllers.ParseBody(receiptUrlRequest, r); err != nil {
		panic(1)
	}

	// Validate request
	// If failed return 422 Unprocessed Entity with error map
	if err := validator.GetValidationErrors(receiptUrlRequest); err != nil {
		controllers.JsonResponse(w, transport.NewValidationError(err), http.StatusUnprocessableEntity)
		return
	}

	// Extract vl parameter
	url, _ := url.Parse(receiptUrlRequest.Url)
	queryParams := url.Query()
	vl := queryParams.Get("vl")

	// Check if receipt with vl param exists in database and if related with authenticated user
	var dbReceipt *models.Receipt
	if dbErr := database.Instance.Where(&models.Receipt{UserID: &user.Id, Vl: vl}).First(&dbReceipt).Error; dbErr == nil {
		// User has already scanned this receipt
		// Return appropriate error response
		controllers.JsonResponse(w, ErrReceiptAlreadyScanned, http.StatusBadRequest)
		return
	}

	// Create receipt with status `pending`
	dbReceipt = &models.Receipt{UserID: &user.Id, Vl: vl, Status: models.RECEIPT_STATUS_PENDING}
	dbResult := database.Instance.Create(dbReceipt)
	if dbResult.Error != nil {
		// Error while creating receipt
		controllers.JsonResponse(w, ErrCannotCreateReceipt, http.StatusBadRequest)
		return
	}

	// Push message Receipt url to SQS
	// And add Receipt ID to message attributes
	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"ReceiptId": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprint((dbReceipt.ID))),
			},
		},
		MessageBody: aws.String(receiptUrlRequest.Url),
		QueueUrl:    pendingReceiptsSqsUrl,
	}

	if _, err := client.SendMessage(sqsMessageInput); err != nil {
		log.Printf("Could not send sqs message: %+v\n", err)
		panic(1)
	}

	// Receipt created successfully
	controllers.JsonInfoResponse(w, "Receipt created and set to be processed", http.StatusCreated)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
