package customer

import (
	"database/sql"
	"net/http"
)

type CustomerModel interface {
	CustomerCreateModel | CustomerUpdateModel
}

func NewCustomer(db *sql.DB) *http.ServeMux {
	return newMux(newHandler(newRepo(db)))
}
