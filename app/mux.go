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

	mux.HandleFunc("GET /{$}", roothandler)
	mux.HandleFunc("POST /api/auth", auth.CreateAuthToken)
	mux.Handle("/", docs.NewMux())
	mux.Handle("/api/customer/", customer.NewCustomer(db))
	mux.Handle("GET /api/secured", auth.Verify(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is secured path"))
	})))

	return mux
}
