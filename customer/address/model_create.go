package address

type AddressCreateModel struct {
	Address    string `json:"address"`
	Address2   string `json:"address2"`
	District   string `json:"district"`
	CityId     int16  `json:"city_id"`
	PostalCode string `json:"postal_code"`
}
