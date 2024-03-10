package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"github.com/steffanturanjanin/receipt-manager/internal/user"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
	validation "github.com/steffanturanjanin/receipt-manager/internal/validator"
)

type RegisterUserRequest struct {
	FirstName string `validate:"required, max=255" json:"first_name"`
	LastName  string `validate:"required, max=255" json:"last_name"`
	Email     string `validate:"required, email, unique=users.email" json:"email"`
	Password  string `validate:"required, min=8, max=100" json:"password"`
}

var (
	// Env variables
	AccessTokenPrivateKey  string
	AccessTokenMaxAge      int
	AccessTokenTTL         time.Duration
	RefreshTokenPrivateKey string
	RefreshTokenMaxAge     int
	RefreshTokenTTL        time.Duration

	// Tokens
	AccessToken  string
	RefreshToken string

	// Auth cookies
	AccessTokenCookie  http.Cookie
	RefreshTokenCookie http.Cookie
	LoggedInCookie     http.Cookie

	// Validator
	Validator *validation.Validator

	// Router
	GorillaLambda *gorillamux.GorillaMuxAdapter

	// Errors
	ErrServiceUnavailable = transport.NewServiceUnavailableError()
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}

	AccessTokenPrivateKey = os.Getenv("AccessTokenPrivateKey")
	AccessTokenMaxAge, _ = strconv.Atoi(os.Getenv("AccessTokenMaxAge"))
	AccessTokenTTL, _ = time.ParseDuration(os.Getenv("AccessTokenExpiresIn"))
	RefreshTokenPrivateKey = os.Getenv("RefreshTokenPrivateKey")
	RefreshTokenMaxAge, _ = strconv.Atoi(os.Getenv("RefreshTokenMaxAge"))
	RefreshTokenTTL, _ = time.ParseDuration(os.Getenv("RefreshTokenExpiresIn"))

	// Initialize Router
	Router := mux.NewRouter()
	Router.HandleFunc("/auth/register", handler)
	GorillaLambda = gorillamux.New(Router)

	// Initialize Validator
	Validator = validation.NewDefaultValidator()
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	registerRequest := &user.RegisterUserRequest{}
	if err := controllers.ParseBody(registerRequest, r); err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Validate request
	// If failed return 422 Unprocessed Entity with error map
	if err := Validator.GetValidationErrors(registerRequest); err != nil {
		controllers.JsonResponse(w, transport.NewValidationError(err), http.StatusUnprocessableEntity)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	user := models.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Email:     registerRequest.Email,
		Password:  hashedPassword,
	}

	// Create user with hashed password
	if err := db.Instance.Create(&user).Error; err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Generate access token
	AccessToken, err = utils.CreateToken(AccessTokenTTL, user.ID, AccessTokenPrivateKey)
	if err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Generate refresh token
	RefreshToken, err = utils.CreateToken(RefreshTokenTTL, user.ID, RefreshTokenPrivateKey)
	if err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	// Generate access token cookie
	AccessTokenCookie = http.Cookie{
		Name:     "access_token",
		Value:    AccessToken,
		Path:     "/",
		MaxAge:   AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	}

	// Generate refresh token cookie
	RefreshTokenCookie = http.Cookie{
		Name:     "refresh_token",
		Value:    RefreshToken,
		Path:     "/",
		MaxAge:   RefreshTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	}

	// Generate logged in cookie
	LoggedInCookie = http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: false,
	}

	// Set auth cookies
	http.SetCookie(w, &AccessTokenCookie)
	http.SetCookie(w, &RefreshTokenCookie)
	http.SetCookie(w, &LoggedInCookie)

	response := dto.AccessToken{AccessToken: AccessToken}

	controllers.JsonResponse(w, response, http.StatusCreated)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *r.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
