package customer

import (
	"database/sql"
	"net/http"
)

type CustomerModel interface {
	CustomerCreateModel | CustomerUpdateModel
}

func NewCustomer(db *sql.DB) *http.ServeMux {
	repo := newRepo(db)
	service := NewService(*repo)
	handler := NewHandler(service)

	return newMux(handler)
}
