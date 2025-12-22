package address

import (
	"errors"
)

type ModelCreate struct {
	Address    *string `json:"address"`
	Address2   *string `json:"address2"`
	District   *string `json:"district"`
	CityId     *int16  `json:"city_id"`
	PostalCode *string `json:"postal_code"`
}

var errAddressAddressMoreThan50Chars = errors.New("address cannot be more than 50 characters")
var errAddressAddressIsNil = errors.New("address is required")
var errAddressAddress2MoreThan50Chars = errors.New("address2 cannot be more than 50 characters")
var errAddressDistrictMoreThan20Chars = errors.New("district name cannot be more than 20 characters")
var errAddressDistrictIsNil = errors.New("district is required")
var errAddressCityIdIsNil = errors.New("city is required")
var errAddressPostalCodeMoreThan10Chars = errors.New("postal code cannot be more than 10 characters")

func (m ModelCreate) Validate() error {
	if m.Address == nil {
		return errAddressAddressIsNil
	}

	if len(*m.Address) > 50 {
		return errAddressAddressMoreThan50Chars
	}

	if m.Address2 != nil && len(*m.Address2) > 50 {
		return errAddressAddress2MoreThan50Chars
	}

	if m.District == nil {
		return errAddressDistrictIsNil
	}

	if len(*m.District) > 20 {
		return errAddressDistrictMoreThan20Chars
	}

	if m.CityId == nil {
		return errAddressCityIdIsNil
	}

	if m.PostalCode != nil && len(*m.PostalCode) > 10 {
		return errAddressPostalCodeMoreThan10Chars
	}

	return nil
}
