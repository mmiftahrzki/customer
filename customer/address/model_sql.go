package address

import "database/sql"

type ModelSQL struct {
	Id         sql.NullInt16
	Address    sql.NullString
	District   sql.NullString
	CityId     sql.NullInt16
	City       sql.NullString
	PostalCode sql.NullString
	CountryId  sql.NullInt16
	Country    sql.NullString
}
