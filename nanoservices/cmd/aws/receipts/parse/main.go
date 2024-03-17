package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/pkg/receipt-fetcher"
)

type ParseReceiptRequest struct {
	Url string `validate:"url" json:"url"`
}

var (
	gorillaLambda *gorillamux.GorillaMuxAdapter
)

func init() {

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	handler := corsMiddleware(jsonMiddleware(handler))

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/receipts/parse", handler)
	gorillaLambda = gorillamux.New(router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	parseReceiptRequest := new(ParseReceiptRequest)

	if err := controllers.ParseBody(parseReceiptRequest, r); err != nil {
		controllers.JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	receipt, err := receipt_fetcher.Get(parseReceiptRequest.Url)
	if err != nil {
		controllers.JsonErrorResponse(w, err)
		return
	}

	controllers.JsonResponse(w, receipt, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
