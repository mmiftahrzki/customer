package address

type AddressUpdateModel struct {
	Address    *string `json:"address"`
	Address2   *string `json:"address2"`
	District   *string `json:"district"`
	PostalCode *string `json:"postal_code"`
}

func ValidateAddressUpdateModel(m AddressUpdateModel) error {
	if m.Address != nil && len(*m.Address) > 50 {
		return errAddressAddressMoreThan50Chars
	}

	if m.Address2 != nil && len(*m.Address2) > 50 {
		return errAddressAddress2MoreThan50Chars
	}

	if m.District != nil && len(*m.District) > 20 {
		return errAddressDistrictMoreThan20Chars
	}

	if m.PostalCode != nil && len(*m.PostalCode) > 10 {
		return errAddressPostalCodeMoreThan10Chars
	}

	return nil
}
