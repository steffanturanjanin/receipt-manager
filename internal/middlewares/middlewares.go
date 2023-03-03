package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

type contextKey string

const (
	CURRENT_USER contextKey = "currentUser"
)

func setAuthMIddleware(next http.HandlerFunc) http.HandlerFunc {
	db := database.Instance
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		var accessToken string
		accessTokenCookie, err := r.Cookie("access_token")

		authorizationHeader := r.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = accessTokenCookie.Value
		}

		if accessToken == "" {
			transport.ResponseJson(w, nil, 401)
		}

		accessTokenPublicKey := os.Getenv("ACCESS_TOKEN_PUBLI_KEY")
		sub, err := utils.ValidateToken(accessToken, accessTokenPublicKey)
		if err != nil {
			transport.ResponseJson(w, nil, 401)
		}

		userId, _ := strconv.Atoi(fmt.Sprint(sub))
		userResponse, err := userService.GetById(userId)
		if err != nil {
			transport.ResponseJson(w, nil, 403)
		}

		ctx := context.WithValue(r.Context(), CURRENT_USER, *userResponse)
		requestWithContext := r.WithContext(ctx)

		next(w, requestWithContext)
	}
}
