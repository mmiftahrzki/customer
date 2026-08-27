package customer

import (
	"database/sql"
)

type customer struct {
	Handler *handler
}

func New(db *sql.DB) customer {
	repo := newRepo(db)
	service := newService(repo)
	handler := newHandler(service)

	return customer{Handler: &handler}
}
