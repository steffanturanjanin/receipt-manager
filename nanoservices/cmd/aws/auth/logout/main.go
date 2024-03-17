package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/steffanturanjanin/receipt-manager/internal/auth"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
)

var (
	authService   *auth.AuthService
	gorillaLambda *gorillamux.GorillaMuxAdapter
)

func init() {
	authService = auth.NewAuthService(nil)

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	authMiddleware := middlewares.SetAuthMiddleware
	handler := authMiddleware(corsMiddleware(jsonMiddleware(handler)))

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/auth/logout", handler)
	gorillaLambda = gorillamux.New(router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	authCookies := authService.Logout()

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.RefreshTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	controllers.JsonResponse(w, nil, http.StatusOK)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *r.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
