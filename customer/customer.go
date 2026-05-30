package customer

import (
	"database/sql"
)

type customer struct {
	Handler handler
}

func New(db *sql.DB) customer {
	return customer{newHandler(newService(newRepo(db)))}
}
