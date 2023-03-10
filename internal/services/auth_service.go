package services

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
)

type AccessTokenCookie http.Cookie
type RefreshTokenCookie http.Cookie
type LoggedInCookie http.Cookie

type AuthCookies struct {
	AccessTokenCookie
	RefreshTokenCookie
	LoggedInCookie
}

type AuthService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewAuthService(userRepository repositories.UserRepositoryInterface) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) RegisterUser(request transport.RegisterUserRequestDTO) (*transport.UserResponseDTO, error) {
	userModel, err := service.UserRepository.Create(request)
	if err != nil {
		return nil, err
	}

	response := models.NewUserResponseDTOFromUserModel(*userModel)

	return &response, nil
}

func (service *AuthService) LoginUser(request transport.LoginUserRequestDTO) (*transport.LoginUserResponseDTO, *AuthCookies, error) {
	userModel, err := service.UserRepository.GetByEmail(strings.ToLower(request.Email))
	if err != nil {
		return nil, nil, err
	}

	if err := utils.VerifyPassword(userModel.Password, request.Password); err != nil {
		return nil, nil, errors.NewErrUnauthorized(err, "Invalid credentials.")
	}

	accessTokenPrivateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenMaxAge, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAX_AGE"))
	accessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRES_IN"))
	if err != nil {
		return nil, nil, err
	}

	accessToken, err := utils.CreateToken(accessTokenTTL, userModel.ID, accessTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenPrivateKey := os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	refreshTokenMaxAge, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAX_AGE"))
	refreshTokenTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRES_IN"))
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := utils.CreateToken(refreshTokenTTL, userModel.ID, refreshTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	accessTokenCookie := AccessTokenCookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   accessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	}

	refreshTokenCookie := RefreshTokenCookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   refreshTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	}

	loggedInCookie := LoggedInCookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   accessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: false,
	}

	return &transport.LoginUserResponseDTO{AccessToken: accessToken}, &AuthCookies{
		AccessTokenCookie:  accessTokenCookie,
		RefreshTokenCookie: refreshTokenCookie,
		LoggedInCookie:     loggedInCookie,
	}, nil
}

func (service *AuthService) RefreshToken(refreshToken RefreshTokenCookie) (*transport.LoginUserResponseDTO, *AuthCookies, error) {
	refreshTokenPublicKey := os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	sub, _ := utils.ValidateToken(refreshToken.Value, refreshTokenPublicKey)
	userId, _ := strconv.Atoi(fmt.Sprint(sub))
	userModel, err := service.UserRepository.GetById(userId)
	if err != nil {
		return nil, nil, err
	}

	accessTokenPrivateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	accessTokenMaxAge, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAX_AGE"))
	accessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRES_IN"))
	if err != nil {
		return nil, nil, err
	}

	accessToken, err := utils.CreateToken(accessTokenTTL, userModel.ID, accessTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	accessTokenCookie := AccessTokenCookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   accessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	}

	loggedInCookie := LoggedInCookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   accessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: false,
	}

	return &transport.LoginUserResponseDTO{AccessToken: accessToken},
		&AuthCookies{
			AccessTokenCookie: accessTokenCookie,
			LoggedInCookie:    loggedInCookie,
		}, nil
}

func (AuthService *AuthService) Logout() AuthCookies {
	accessTokenCookie := AccessTokenCookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}

	refreshTokenCookie := RefreshTokenCookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}

	loggedInCookie := LoggedInCookie{
		Name:     "logged_in",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: false,
	}

	return AuthCookies{
		AccessTokenCookie:  accessTokenCookie,
		RefreshTokenCookie: refreshTokenCookie,
		LoggedInCookie:     loggedInCookie,
	}
}
