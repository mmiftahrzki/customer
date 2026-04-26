package auth

import (
	"context"
	"net/http"
	"slices"

	pkgMiddleware "github.com/mmiftahrzki/customer/middleware"
	"github.com/mmiftahrzki/customer/responses"
)

type middleware struct {
	service service
}

func newMiddleware(svc service) middleware {
	return middleware{
		service: svc,
	}
}

func (m *middleware) VerifyJWT(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authValue := r.Header.Get(RequestHeaderAuthKey)
		tokenStr, err := extractAuthTokenStr(authValue)
		if err != nil {
			responses.Error(w, http.StatusUnauthorized, err.Error())

			return
		}

		token, err := m.service.getToken(tokenStr)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err.Error())

			return
		}

		r = r.WithContext(context.WithValue(r.Context(), JWTContextKey, token.Claims))

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) RequiredRole(roles ...string) pkgMiddleware.Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			jwtCtx := r.Context().Value(JWTContextKey)
			claim, ok := jwtCtx.(*ModelClaim)
			if !ok {
				responses.Error(w, http.StatusUnauthorized, "invalid jwt claims")
				return
			}

			if !slices.Contains(roles, claim.Role) {
				responses.Error(w, http.StatusForbidden, "forbidden")

				return
			}

			hf.ServeHTTP(w, r)
		}
	}
}
