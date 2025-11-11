package customer

import "database/sql"

type customer struct {
	Handler handler
}

func New(db *sql.DB) customer {
	return customer{Handler: newHandler(newService(newRepo(db)))}
}
