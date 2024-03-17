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

	"github.com/steffanturanjanin/receipt-manager/internal/auth"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/user"
	validation "github.com/steffanturanjanin/receipt-manager/internal/validator"
)

var (
	// Infrastructure
	userRepository *user.UserRepository
	authService    *auth.AuthService

	// Validator
	validator *validation.Validator

	// Router
	gorillaLambda *gorillamux.GorillaMuxAdapter
)

func init() {
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	userRepository = user.NewUserRepository(db.Instance)
	authService = auth.NewAuthService(userRepository)
	validator = validation.NewDefaultValidator()

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	handler := corsMiddleware(jsonMiddleware(handler))

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/auth/login", handler).Methods("POST")
	gorillaLambda = gorillamux.New(router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	loginRequest := new(user.LoginUserRequest)

	if err := controllers.ParseBody(loginRequest, r); err != nil {
		controllers.JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	if err := controllers.ValidateRequest(loginRequest, validator); err != nil {
		controllers.JsonErrorResponse(w, err)
		return
	}

	response, authCookies, err := authService.LoginUser(*loginRequest)
	if err != nil {
		controllers.JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.RefreshTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	controllers.JsonResponse(w, response, http.StatusCreated)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *r.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
