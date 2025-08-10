package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mmiftahrzki/go-rest-api/response"
)

type jwtContextKey int

const key jwtContextKey = iota
const request_header_auth_key string = "Authorization"

var errEmptyAuth = errors.New("authorization header not found")
var errInvalidAuth = errors.New("invalid authorization header")

func extractAuthTokenStr(auth_value string) (string, error) {
	var token_str string

	if len(auth_value) == 0 {
		return token_str, errEmptyAuth
	}

	auth_value_fields := strings.Fields(auth_value)
	if len(auth_value_fields) != 2 || auth_value_fields[0] != "Bearer" {
		return token_str, errInvalidAuth
	}

	token_str = auth_value_fields[1]

	return token_str, nil
}

func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := response.New()

		auth_value := r.Header.Get(request_header_auth_key)
		token_str, err := extractAuthTokenStr(auth_value)
		if err != nil {
			response.Message = err.Error()

			w.WriteHeader(http.StatusBadRequest)
			w.Write(response.ToJson())

			return
		}

		token, err := jwt.ParseWithClaims(token_str, &AuthClaimModel{}, func(t *jwt.Token) (any, error) {
			method, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method")
			}

			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			response.Message = err.Error()

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(response.ToJson()))

			return
		}

		if !token.Valid {
			response.Message = "invalid jwt"

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(response.ToJson()))

			return
		}

		claims, ok := token.Claims.(*AuthClaimModel)
		if !ok {
			response.Message = "invalid jwt claims"

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(response.ToJson()))

			return
		}

		r = r.WithContext(context.WithValue(r.Context(), key, claims))
		next.ServeHTTP(w, r)
	})
}
