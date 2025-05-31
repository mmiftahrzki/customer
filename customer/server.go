package customer

import (
	"database/sql"
	"net/http"
)

type server struct {
	Mux *http.ServeMux
}

func New(db *sql.DB) *server {
	server := &server{
		Mux: newMux(newHandler(newRepo(db))),
	}

	return server
}
