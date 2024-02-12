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
	"gorm.io/gorm"

	"github.com/steffanturanjanin/receipt-manager/internal/auth"
	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/user"
	validation "github.com/steffanturanjanin/receipt-manager/internal/validator"
)

var (
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost     = os.Getenv("DB_HOST")
	dbPort     = os.Getenv("DB_PORT")
	dbName     = os.Getenv("DB_NAME")

	db             *gorm.DB
	err            error
	userRepository *user.UserRepository
	authService    *auth.AuthService
	validator      *validation.Validator
	gorillaLambda  *gorillamux.GorillaMuxAdapter
)

func init() {
	db, err = database.InitDB(dbName, dbUser, dbPassword, dbHost, dbPort)
	if err != nil {
		os.Exit(1)
	}

	userRepository = user.NewUserRepository(db)
	authService = auth.NewAuthService(userRepository)
	validator = validation.NewDefaultValidator()

	r := mux.NewRouter()

	r.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		registerRequest := new(user.RegisterUserRequest)

		if err := controllers.ParseBody(registerRequest, r); err != nil {
			controllers.JsonErrorResponse(w, errors.NewHttpError(err))
			return
		}

		if err := controllers.ValidateRequest(registerRequest, validator); err != nil {
			controllers.JsonErrorResponse(w, err)
			return
		}

		response, err := authService.RegisterUser(*registerRequest)
		if err != nil {
			controllers.JsonErrorResponse(w, errors.NewHttpError(err))
			return
		}

		controllers.JsonResponse(w, response, http.StatusCreated)
	})

	gorillaLambda = gorillamux.New(r)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *r.Version1(), err
}

func main() {
	lambda.Start(Handler)
}