package customer

import (
	"database/sql"

	"github.com/mmiftahrzki/customer/customer/address"
)

type sqlModel struct {
	id        sql.NullInt16
	firstName sql.NullString
	lastName  sql.NullString
	email     sql.NullString
	addressId sql.NullInt16
	address   address.SQLModel
	active    sql.NullBool
	createdAt sql.NullTime
}
