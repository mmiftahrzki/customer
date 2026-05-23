package address

type ModelRead struct {
	Id         int    `json:"id"`
	Address    string `json:"address"`
	Address2   string `json:"address_2"`
	District   string `json:"district"`
	CityId     int    `json:"city_id"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	CountryId  int    `json:"country_id"`
	Country    string `json:"country"`
}
