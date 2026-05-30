package customer

import (
	"database/sql"

	"github.com/mmiftahrzki/customer/customer/address"
)

type modelSQL struct {
	id        sql.NullInt16
	firstName sql.NullString
	lastName  sql.NullString
	email     sql.NullString
	addressId sql.NullInt16
	address   address.ModelSQL
	active    sql.NullBool
	createdAt sql.NullTime
}
