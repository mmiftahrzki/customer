package app

import (
	"database/sql"
	"net/http"

	"github.com/mmiftahrzki/customer/auth"
	"github.com/mmiftahrzki/customer/customer"
	"github.com/mmiftahrzki/customer/docs"
)

func newMux(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	auth := auth.New()
	customer := customer.New(db, auth.Middleware.VerifyJWT)

	mux.HandleFunc("GET /{$}", roothandler)
	mux.Handle("/", docs.NewMux())
	mux.Handle("/api/auth/", auth.Handler)
	mux.Handle("/api/customer/", customer.Handler)

	return mux
}
