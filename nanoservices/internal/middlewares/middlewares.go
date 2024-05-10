package middlewares

import (
	"context"
	native_errors "errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/steffanturanjanin/receipt-manager/internal/controllers"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/query"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
)

type ContextKey string

const (
	CURRENT_USER     ContextKey = "CURRENT_USER"
	PAGINATION_QUERY ContextKey = "PAGINATION_QUERY"
	SORT_QUERY       ContextKey = "SORT_QUERY"
)

func SetJsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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
			controllers.JsonErrorResponse(w, errors.NewErrUnauthorized(nil, "No Access Token."))
			return
		}

		accessTokenPublicKey := os.Getenv("AccessTokenPublicKey")
		sub, err := utils.ValidateToken(accessToken, accessTokenPublicKey)
		if err != nil {
			controllers.JsonErrorResponse(w, errors.NewErrUnauthorized(err, "Invalid access token."))
			return
		}

		userId, _ := strconv.Atoi(fmt.Sprint(sub))
		userResponse, err := userService.GetById(userId)
		if err != nil {
			var appError errors.AppErrorInterface
			if native_errors.As(err, &appError) {
				controllers.JsonErrorResponse(
					w,
					errors.NewHttpError(errors.NewErrUnauthorized(appError.GetError(), "Forbidden.")),
				)
			}

			controllers.JsonErrorResponse(w, errors.NewErrForbidden(nil, "Forbidden"))
			return
		}

		ctx := context.WithValue(r.Context(), CURRENT_USER, *userResponse)
		requestWithContext := r.WithContext(ctx)

		next(w, requestWithContext)
	}
}

func SetQueryParamsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortQuery := query.SortQueryFromRequest(*r)
		paginationQuery := query.PaginationQueryFromRequest(*r)

		var ctx context.Context
		ctx = context.WithValue(r.Context(), PAGINATION_QUERY, paginationQuery)
		ctx = context.WithValue(ctx, SORT_QUERY, sortQuery)

		next(w, r.WithContext(ctx))
	}
}

func SetCorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next(w, r)
	}
}

func GetAuthUser(r *http.Request) dto.User {
	return r.Context().Value(CURRENT_USER).(dto.User)
}

func GetPaginationQueryParams(r *http.Request) query.PaginationQuery {
	return r.Context().Value(PAGINATION_QUERY).(query.PaginationQuery)
}

func GetSortQueryParams(r *http.Request) *query.SortQuery {
	return r.Context().Value(SORT_QUERY).(*query.SortQuery)
}
