package customer

import "database/sql"

type CustomerSql struct {
	customer_id sql.NullInt16
	first_name  sql.NullString
	last_name   sql.NullString
	email       sql.NullString
	address_id  sql.NullInt16
	active      sql.NullBool
	create_date sql.NullTime
}
