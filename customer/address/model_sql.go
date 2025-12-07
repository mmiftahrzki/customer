package address

import "database/sql"

type SQLModel struct {
	Id         sql.NullInt16
	Address    sql.NullString
	District   sql.NullString
	CityId     sql.NullInt16
	PostalCode sql.NullString
}
