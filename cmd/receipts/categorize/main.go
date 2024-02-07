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
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
)

var (
	categorizeService services.CategorizeService
	gorillaLambda     *gorillamux.GorillaMuxAdapter
)

func init() {
	err := database.InitializeDB()
	if err != nil {
		panic(1)
	}

	router := mux.NewRouter()
	router.HandleFunc("/receipts/categorize", handler2).Methods("POST")
	gorillaLambda = gorillamux.New(router)

	categorizeService = services.CategorizeService{}
}

type RequestUrl struct {
	Url string `json:"url"`
}

// func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	var requestUrl RequestUrl
// 	err := json.Unmarshal([]byte(request.Body), &requestUrl)
// 	if err != nil {
// 		panic(1)
// 	}

// 	requestUrlJson, err := json.Marshal(requestUrl)
// 	if err != nil {
// 		panic(1)
// 	}

// 	response := events.APIGatewayProxyResponse{
// 		StatusCode:        200,
// 		Headers:           make(map[string]string, 0),
// 		MultiValueHeaders: make(map[string][]string, 0),
// 		Body:              Test,
// 		IsBase64Encoded:   false,
// 	}

// 	return response, nil

// 	receipt, err := receipt_fetcher.Get(requestUrl.Url)
// 	if err != nil {
// 		panic(1)
// 	}

// 	receiptItemInputs := make([]services.ReceiptItemIn, 0)
// 	for _, receiptItem := range receipt.Items {
// 		receiptItemInputs = append(receiptItemInputs, services.ReceiptItemIn{Name: receiptItem.Name})
// 	}

// 	var categories []models.Category
// 	dbResult := database.Instance.Find(&categories)
// 	if dbResult.Error != nil {
// 		log.Printf("Error while reading categories %+v\n", dbResult.Error)
// 		return response, dbResult.Error
// 	}

// 	categoryInputs := make([]services.CategoryIn, 0)
// 	for _, category := range categories {
// 		categoryInputs = append(categoryInputs, services.CategoryIn{Id: category.ID, Name: category.Name})
// 	}

// 	categorizedMap, err := categorizeService.Categorize(receiptItemInputs, categoryInputs)
// 	if err != nil {
// 		log.Printf("Error while trying to categorize items %+v\n", err)
// 		return response, err
// 	}

// 	log.Printf("%+v\n", categorizedMap)

// 	return response, nil
// }

func handler2(w http.ResponseWriter, r *http.Request) {
	var req RequestUrl
	if err := controllers.ParseBody(&req, r); err != nil {
		panic(1)
	}

	controllers.JsonResponse(w, req, 200)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
