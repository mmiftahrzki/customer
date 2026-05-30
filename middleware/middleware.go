package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Add(middleware Middleware, hf http.HandlerFunc) http.HandlerFunc {
	return middleware(hf)
}

func Pipe(middlewares ...Middleware) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			hf = middlewares[i](hf)
		}

		return hf
	}
}
