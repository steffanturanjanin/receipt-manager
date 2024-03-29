package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/user"
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
	UserRepository user.UserRepositoryInterface
}

func NewAuthService(userRepository user.UserRepositoryInterface) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) RegisterUser(request user.RegisterUserRequest) (*user.UserDto, error) {
	userModel, err := service.UserRepository.Create(request)
	if err != nil {
		return nil, err
	}

	response := user.NewUserResponseDTOFromUserModel(*userModel)

	return &response, nil
}

func (service *AuthService) LoginUser(request user.LoginUserRequest) (*user.AccessToken, *AuthCookies, error) {
	userModel, err := service.UserRepository.GetByEmail(strings.ToLower(request.Email))
	if err != nil {
		return nil, nil, err
	}

	if err := utils.VerifyPassword(userModel.Password, request.Password); err != nil {
		return nil, nil, errors.NewErrUnauthorized(err, "Invalid credentials.")
	}

	accessTokenPrivateKey := os.Getenv("AccessTokenPrivateKey")
	accessTokenMaxAge, _ := strconv.Atoi(os.Getenv("AccessTokenMaxAge"))
	accessTokenTTL, err := time.ParseDuration(os.Getenv("AccessTokenExpiresIn"))
	if err != nil {
		return nil, nil, err
	}

	accessToken, err := utils.CreateToken(accessTokenTTL, userModel.ID, accessTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenPrivateKey := os.Getenv("RefreshTokenPrivateKey")
	refreshTokenMaxAge, _ := strconv.Atoi(os.Getenv("RefreshTokenMaxAge"))
	refreshTokenTTL, err := time.ParseDuration(os.Getenv("RefreshTokenExpiresIn"))
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

	return &user.AccessToken{AccessToken: accessToken}, &AuthCookies{
		AccessTokenCookie:  accessTokenCookie,
		RefreshTokenCookie: refreshTokenCookie,
		LoggedInCookie:     loggedInCookie,
	}, nil
}

func (service *AuthService) RefreshToken(refreshToken RefreshTokenCookie) (*user.AccessToken, *AuthCookies, error) {
	refreshTokenPublicKey := os.Getenv("RefreshTokenPublicKey")
	sub, _ := utils.ValidateToken(refreshToken.Value, refreshTokenPublicKey)
	userId, _ := strconv.Atoi(fmt.Sprint(sub))
	userModel, err := service.UserRepository.GetById(userId)
	if err != nil {
		return nil, nil, err
	}

	accessTokenPrivateKey := os.Getenv("AccessTokenPrivateKey")
	accessTokenMaxAge, _ := strconv.Atoi(os.Getenv("AccessTokenMaxAge"))
	accessTokenTTL, err := time.ParseDuration(os.Getenv("AccessTokenExpiresIn"))
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

	return &user.AccessToken{AccessToken: accessToken},
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
