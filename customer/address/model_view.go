package address

type Address struct {
	Id         uint8  `json:"id"`
	Address    string `json:"address"`
	Address2   string `json:"address_2"`
	District   string `json:"district"`
	CityId     uint8  `json:"city_id"`
	PostalCode string `json:"postal_code"`
}
