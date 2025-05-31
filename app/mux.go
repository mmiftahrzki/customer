package app

import (
	"database/sql"
	"net/http"

	"github.com/mmiftahrzki/customer/customer"
	"github.com/mmiftahrzki/customer/docs"
)

func newMux(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", roothandler)
	mux.Handle("/api/customer/", customer.New(db).Mux)
	mux.Handle("/", docs.NewMux())

	return mux
}
