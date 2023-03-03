package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

type AuthController struct {
	AuthService *services.AuthService
	Validator   *validator.Validate
}

func NewAuthController(authService *services.AuthService, validator *validator.Validate) *AuthController {
	return &AuthController{
		AuthService: authService,
		Validator:   validator,
	}
}

func (controller *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerRequest transport.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		transport.ResponseJson(w, transport.HttpErrorResponse{
			Message: "Could not parse the Register User Request.",
			Code:    http.StatusInternalServerError,
		}, http.StatusInternalServerError)
	}

	err := controller.Validator.Struct(registerRequest)
	if err != nil {
		transport.ValidationErrorResponseJson(w, err)
		return
	}

	response, err := controller.AuthService.RegisterUser(registerRequest)
	if err != nil {
		transport.ResponseJson(w, transport.HttpErrorResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, http.StatusInternalServerError)
	}

	transport.ResponseJson(w, response, http.StatusCreated)
}

func (controller *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest transport.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		transport.ResponseJson(w, transport.HttpErrorResponse{
			Message: "Could not parse the Login Request.",
			Code:    http.StatusInternalServerError,
		}, http.StatusInternalServerError)
	}

	err := controller.Validator.Struct(loginRequest)
	if err != nil {
		transport.ValidationErrorResponseJson(w, err)
		return
	}

	loginResponse, authCookies, err := controller.AuthService.LoginUser(loginRequest)
	if err != nil {
		transport.ResponseJson(w, transport.HttpErrorResponse{
			Message: "Failed authenticating user",
			Code:    http.StatusBadRequest,
		}, http.StatusBadRequest)
	}

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.RefreshTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	transport.ResponseJson(w, loginResponse, http.StatusOK)
}

func (controller *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, _ := r.Cookie("refresh_token")

	loginResponse, authCookies, err := controller.AuthService.RefreshToken(services.RefreshTokenCookie(*refreshTokenCookie))
	if err != nil {
		transport.ResponseJson(w, transport.HttpErrorResponse{
			Message: "Failed refreshing access token",
			Code:    http.StatusBadRequest,
		}, http.StatusBadRequest)
	}

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	transport.ResponseJson(w, loginResponse, http.StatusOK)
}

func (controller *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	authCookies := controller.AuthService.Logout()

	http.SetCookie(w, (*http.Cookie)(&authCookies.AccessTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.RefreshTokenCookie))
	http.SetCookie(w, (*http.Cookie)(&authCookies.LoggedInCookie))

	transport.ResponseJson(w, nil, http.StatusOK)
}
