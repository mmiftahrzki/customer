package app

import (
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Test(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware test: it works")

		next.ServeHTTP(w, r)
	}
}

func add(middleware Middleware, hf http.HandlerFunc) http.HandlerFunc {
	return middleware(hf)
}

func pipe(middlewares ...Middleware) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			hf = middlewares[i](hf)
		}

		return hf
	}
}
