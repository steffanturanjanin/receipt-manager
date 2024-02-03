package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	appErrors "github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	receiptFetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

type ReceiptUrlRequest struct {
	Url string `validate:"url" json:"url"`
}

const (
	RECEIPT_URLS_QUEUE   = "receipt_urls"
	RECEIPT_PARSED_QUEUE = "receipt_parsed"
)

var (
	session       *awsSession.Session
	client        *sqs.SQS
	gorillaLambda *gorillamux.GorillaMuxAdapter

	receiptUrlQueueUrl    *string
	receiptParsedQueueUrl *string
)

func init() {
	err := database.InitializeDB()
	if err != nil {
		os.Exit(1)
	}

	router := mux.NewRouter()
	router.HandleFunc("/receipts/url", middlewares.SetAuthMiddleware(handler)).Methods("POST")
	gorillaLambda = gorillamux.New(router)

	sessionOptions := awsSession.Options{SharedConfigState: awsSession.SharedConfigDisable}
	session = awsSession.Must(awsSession.NewSessionWithOptions(sessionOptions))
	client = sqs.New(session)

	if urlResult, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RECEIPT_URLS_QUEUE),
	}); err != nil {
		panic(1)
	} else {
		receiptUrlQueueUrl = urlResult.QueueUrl
	}

	if urlResult, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RECEIPT_PARSED_QUEUE),
	}); err != nil {
		panic(1)
	} else {
		receiptParsedQueueUrl = urlResult.QueueUrl
	}
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middlewares.CURRENT_USER).(dto.User)

	receiptUrlRequest := &ReceiptUrlRequest{}
	if err := controllers.ParseBody(receiptUrlRequest, r); err != nil {
		controllers.JsonErrorResponse(w, appErrors.NewHttpError(err))
		return
	}

	receipt, err := receiptFetcher.Get(receiptUrlRequest.Url)
	if err != nil {
		// If url is invalid
		if errors.Is(err, receiptFetcher.ErrInvalidReceiptUrl) {
			clientError := appErrors.NewErrBadRequest(err, "invalid receipt url")
			controllers.JsonErrorResponse(w, appErrors.NewHttpError(clientError))
			return
		}

		// If receipt data is not available at the moment
		// push message to `receipt_url` queue where it will periodically be checked
		// and eventually processed or rejected
		if errors.Is(err, receiptFetcher.ErrReceiptDataNotAvailable) {
			sqsMessageInput := &sqs.SendMessageInput{
				DelaySeconds: aws.Int64(0),
				MessageAttributes: map[string]*sqs.MessageAttributeValue{
					"UserId": {
						DataType:    aws.String("Number"),
						StringValue: aws.String(fmt.Sprint(user.Id)),
					},
				},
				MessageBody: aws.String(receiptUrlRequest.Url),
				QueueUrl:    receiptUrlQueueUrl,
			}

			_, err := client.SendMessage(sqsMessageInput)
			if err != nil {
				controllers.JsonErrorResponse(w, appErrors.NewHttpError(err))
				return
			}
		}
	}

	// Check if user has scanned this receipt before
	var dbReceipt *models.Receipt
	database.Instance.Where(&models.Receipt{UserID: user.Id, PfrNumber: receipt.Number}).First(&dbReceipt)
	if dbReceipt != nil {
		controllers.JsonInfoResponse(w, "Receipt has already been scanned", http.StatusUnprocessableEntity)
		return
	}

	receiptJsonSerialized, err := json.Marshal(receipt)
	if err != nil {
		panic(1)
	}

	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprint(user.Id)),
			},
		},
		MessageBody: aws.String(string(receiptJsonSerialized)),
		QueueUrl:    receiptParsedQueueUrl,
	}

	_, err = client.SendMessage(sqsMessageInput)
	if err != nil {
		controllers.JsonErrorResponse(w, err)
		return
	}

	controllers.JsonInfoResponse(w, "Receipt is set to be processed.", http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
