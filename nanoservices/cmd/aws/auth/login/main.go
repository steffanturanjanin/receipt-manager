package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/middlewares"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
	validation "github.com/steffanturanjanin/receipt-manager/internal/validator"
)

type LoginRequest struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

const (
	INVALID_CREDENTIALS = "Invalid credentials"
)

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
	ErrInvalidCredentials = transport.NewValidationError(map[string]string{
		"email":    INVALID_CREDENTIALS,
		"password": INVALID_CREDENTIALS,
	})
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

	// Initialize validator
	Validator = validation.NewDefaultValidator()

	// Build middleware chain
	jsonMiddleware := middlewares.SetJsonMiddleware
	corsMiddleware := middlewares.SetCorsMiddleware
	handler := corsMiddleware(jsonMiddleware(handler))

	// Initialize router
	Router := mux.NewRouter()
	Router.HandleFunc("/auth/login", handler).Methods("POST")
	GorillaLambda = gorillamux.New(Router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	loginRequest := &LoginRequest{}

	if err := controllers.ParseBody(loginRequest, r); err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusServiceUnavailable)
		return
	}

	if err := controllers.ValidateRequest(loginRequest, Validator); err != nil {
		controllers.JsonErrorResponse(w, err)
		return
	}

	user := &models.User{}

	err := db.Instance.Where("email = ?", loginRequest.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		controllers.JsonResponse(w, ErrInvalidCredentials, http.StatusUnprocessableEntity)
		return
	}

	err = utils.VerifyPassword(user.Password, loginRequest.Password)
	if err != nil {
		controllers.JsonResponse(w, ErrInvalidCredentials, http.StatusUnprocessableEntity)
		return
	}

	accessToken, err := utils.CreateToken(AccessTokenTTL, user.ID, AccessTokenPrivateKey)
	if err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusUnprocessableEntity)
		return
	}

	refreshToken, err := utils.CreateToken(AccessTokenTTL, user.ID, AccessTokenPrivateKey)
	if err != nil {
		controllers.JsonResponse(w, ErrServiceUnavailable, http.StatusUnprocessableEntity)
		return
	}

	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   RefreshTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	}

	loggedInCookie := http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: false,
	}

	http.SetCookie(w, (*http.Cookie)(&accessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&refreshTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&loggedInCookie))

	response := dto.AccessToken{AccessToken: accessToken}

	controllers.JsonResponse(w, response, http.StatusCreated)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := GorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&request))
	return *r.Version1(), err
}

func main() {
	lambda.Start(Handler)
}
