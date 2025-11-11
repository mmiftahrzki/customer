package app

import (
	"database/sql"
	"net/http"

	"github.com/mmiftahrzki/customer/auth"
	"github.com/mmiftahrzki/customer/customer"
	"github.com/mmiftahrzki/customer/docs"
	"github.com/mmiftahrzki/customer/middleware"
)

func chainMiddleware(handler http.HandlerFunc, middlewares ...middleware.Middleware) http.HandlerFunc {
	var handlers http.HandlerFunc = handler

	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]

		handlers = middleware(handlers)
	}

	return handlers
}

func newMux(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	auth := auth.New()
	customer := customer.New(db)

	mux.HandleFunc("GET /{$}", roothandler)
	mux.Handle("/", docs.NewMux())

	mux.HandleFunc("/api/auth", auth.Handler.CreateAuthToken)

	deleteSingleById := chainMiddleware(customer.Handler.DeleteSingleById, auth.Middleware.VerifyJWT)
	mux.HandleFunc("DELETE /api/customer/{id}", deleteSingleById)
	mux.HandleFunc("POST /api/customer/", customer.Handler.PostSingle)
	mux.HandleFunc("GET /api/customer/", customer.Handler.GetMultiple)
	mux.HandleFunc("GET /api/customer/{id}", customer.Handler.GetSingleById)
	mux.HandleFunc("GET /api/customer/{id}/prev", customer.Handler.GetMultiplePrev)
	mux.HandleFunc("GET /api/customer/{id}/next", customer.Handler.GetMultipleNext)
	mux.HandleFunc("PUT /api/customer/{id}", customer.Handler.PutSingleById)
	mux.HandleFunc("PATCH /api/customer/{customer_id}/address/{address_id}", customer.Handler.GetSingleAndUpdateAddressById)

	return mux
}
