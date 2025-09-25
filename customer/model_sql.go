package customer

import (
	"database/sql"

	"github.com/mmiftahrzki/customer/customer/address"
)

type customerSQLModel struct {
	customer_id sql.NullInt16
	first_name  sql.NullString
	last_name   sql.NullString
	email       sql.NullString
	address_id  sql.NullInt16
	address     address.AddressSqlModel
	active      sql.NullBool
	create_date sql.NullTime
}
