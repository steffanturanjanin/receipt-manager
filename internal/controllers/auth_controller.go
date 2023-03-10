package controllers

import (
	"net/http"

	transport "github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/validator"
)

type AuthController struct {
	AuthService *services.AuthService
	Validator   *validator.Validator
}

func NewAuthController(authService *services.AuthService, validator *validator.Validator) *AuthController {
	return &AuthController{
		AuthService: authService,
		Validator:   validator,
	}
}

func (controller *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	registerRequest := new(transport.RegisterUserRequestDTO)

	if err := ParseBody(registerRequest, r); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	if err := ValidateRequest(registerRequest, controller.Validator); err != nil {
		JsonErrorResponse(w, err)
		return
	}

	response, err := controller.AuthService.RegisterUser(*registerRequest)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	JsonResponse(w, response, http.StatusCreated)
}

func (controller *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	loginRequest := new(transport.LoginUserRequestDTO)

	if err := ParseBody(loginRequest, r); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))

		return
	}

	if err := ValidateRequest(loginRequest, controller.Validator); err != nil {
		JsonErrorResponse(w, err)

		return
	}

	loginResponse, authCookies, err := controller.AuthService.LoginUser(*loginRequest)
	if err != nil {
		JsonErrorResponse(w, err)

		return
	}

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.RefreshTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	JsonResponse(w, loginResponse, http.StatusOK)
}

func (controller *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, _ := r.Cookie("refresh_token")

	loginResponse, authCookies, err := controller.AuthService.RefreshToken(services.RefreshTokenCookie(*refreshTokenCookie))
	if err != nil {
		// transport.ResponseJson(w, transport.HttpErrorResponse{
		// 	Message: "Failed refreshing access token",
		// 	Code:    http.StatusBadRequest,
		// }, http.StatusBadRequest)
		JsonErrorResponse(w, errors.NewHttpError(err))

		return
	}

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	JsonResponse(w, loginResponse, http.StatusOK)
}

func (controller *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	authCookies := controller.AuthService.Logout()

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.RefreshTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	JsonResponse(w, nil, http.StatusOK)
}
