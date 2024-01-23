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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	app_errors "github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

type ReceiptUrlRequest struct {
	Url string `validate:"url" json:"url"`
}

const (
	RECEIPT_URLS_QUEUE   = "receipt_urls"
	RECEIPT_PARSED_QUEUE = "receipt_parsed"
)

var (
	Session       *session.Session
	SqsService    *sqs.SQS
	gorillaLambda *gorillamux.GorillaMuxAdapter
)

func init() {
	err := database.InitializeDB()
	if err != nil {
		os.Exit(1)
	}

	SessionOptions := session.Options{SharedConfigState: session.SharedConfigDisable}
	Session = session.Must(session.NewSessionWithOptions(SessionOptions))
	SqsService = sqs.New(Session)

	router := mux.NewRouter()
	router.HandleFunc("/receipts/url", middlewares.SetAuthMiddleware(handler)).Methods("POST")

	gorillaLambda = gorillamux.New(router)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middlewares.CURRENT_USER).(dto.User)

	parseReceiptRequest := new(ReceiptUrlRequest)
	if err := controllers.ParseBody(parseReceiptRequest, r); err != nil {
		controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
		return
	}

	receipt, err := receipt_fetcher.Get(parseReceiptRequest.Url)
	if err != nil {
		// If url is invalid return error response
		if errors.Is(err, receipt_fetcher.ErrInvalidReceiptUrl) {
			errBadRequest := app_errors.NewErrBadRequest(err, err.Error())
			controllers.JsonErrorResponse(w, app_errors.NewHttpError(errBadRequest))
			return
		}

		// If receipt is not available at the moment, push it to the pending_receipts queue
		// Message should contain current user's id and receipt url
		if errors.Is(err, receipt_fetcher.ErrReceiptDataNotAvailable) {
			urlResult, err := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
				QueueName: aws.String(RECEIPT_URLS_QUEUE),
			})

			if err != nil {
				controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
				return
			}

			sqsMessageInput := &sqs.SendMessageInput{
				DelaySeconds: aws.Int64(0),
				MessageAttributes: map[string]*sqs.MessageAttributeValue{
					"UserId": {
						DataType:    aws.String("Number"),
						StringValue: aws.String(fmt.Sprint(user.Id)),
					},
				},
				MessageBody: aws.String(parseReceiptRequest.Url),
				QueueUrl:    urlResult.QueueUrl,
			}

			_, err = SqsService.SendMessage(sqsMessageInput)
			if err != nil {
				controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
				return
			}

			controllers.JsonResponse(w, "Receipt data currently unavailable. It will be processed as soon as possible.", http.StatusOK)
		}

		controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
		return
	}

	urlResult, err := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RECEIPT_PARSED_QUEUE),
	})

	if err != nil {
		controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
		return
	}

	receiptParsedMsg, err := json.Marshal(receipt)
	if err != nil {
		controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
		return
	}

	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprint(user.Id)),
			},
		},
		MessageBody: aws.String(string(receiptParsedMsg)),
		QueueUrl:    urlResult.QueueUrl,
	}

	_, err = SqsService.SendMessage(sqsMessageInput)
	if err != nil {
		controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
		return
	}

	controllers.JsonResponse(w, "Receipt is set to be processed.", http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
