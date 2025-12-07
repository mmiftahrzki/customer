package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	var handlers http.HandlerFunc = handler

	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]

		handlers = middleware(handlers)
	}

	return handlers
}
