package auth

import (
	"net/http"
)

type auth struct {
	http.Handler
	Middleware middleware
}

func New() auth {
	service := newService()
	handler := newHandler(service)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth", handler.CreateAuthToken)

	return auth{
		Middleware: newMiddleware(service),
		Handler:    mux,
	}
}
