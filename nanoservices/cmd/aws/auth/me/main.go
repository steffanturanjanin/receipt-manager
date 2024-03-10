package main

import (
	"context"
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
)

var (
	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/auth/me", middlewares.SetAuthMiddleware(handler)).Methods("GET")
	GorillaLambda = gorillamux.New(Router)
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	// Get Auth user
	user := middlewares.GetAuthUser(r)

	// Return user response
	controllers.JsonResponse(w, user, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *response.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
