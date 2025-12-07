package auth

import (
	"context"
	"net/http"

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
			responses.Error(w, http.StatusBadRequest, err.Error())

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
