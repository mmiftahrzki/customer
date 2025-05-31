package address

import "database/sql"

type AddressSql struct {
	AddressId  sql.NullInt16
	Address    sql.NullString
	Address2   sql.NullString
	District   sql.NullString
	CityId     sql.NullInt16
	PostalCode sql.NullString
}
