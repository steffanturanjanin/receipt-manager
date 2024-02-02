package main

import (
	"context"
	"fmt"
	"log"
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
)

type ReceiptUrlRequest struct {
	Url string `validate:"url" json:"url"`
}

const (
	RECEIPT_URLS_QUEUE = "receipt_urls"
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
	log.Printf("User: %+v\n", user)

	receiptUrlRequest := &ReceiptUrlRequest{}
	if err := controllers.ParseBody(receiptUrlRequest, r); err != nil {
		controllers.JsonErrorResponse(w, app_errors.NewHttpError(err))
		return
	}

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
		MessageBody: aws.String(receiptUrlRequest.Url),
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
