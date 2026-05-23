package address

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ModelUpdate struct {
	Address    *string `json:"address"`
	District   *string `json:"district"`
	PostalCode *string `json:"postal_code"`
	CityId     *uint16 `json:"city_id"`
}

func (this *ModelUpdate) UnmarshalJSON(data []byte) error {
	type Alias ModelUpdate

	var model Alias

	err := json.Unmarshal(data, &model)
	if err != nil {
		return err
	}

	if model.Address == nil {
		return errors.New("address is nil")
	}

	if model.District == nil {
		return errors.New("district is nil")
	}

	if model.PostalCode == nil {
		return errors.New("postal code is nil")
	}

	if model.CityId == nil {
		return errors.New("city is nil")
	}

	*this = ModelUpdate(model)

	err = this.Validate()
	if err != nil {
		return fmt.Errorf("address detail is invalid: %w", err)
	}

	return nil
}

func (this *ModelUpdate) Validate() error {
	if len(*this.Address) > 50 {
		return errAddressAddressMoreThan50Chars
	}

	if len(*this.District) > 20 {
		return errAddressDistrictMoreThan20Chars
	}

	if len(*this.PostalCode) > 10 {
		return errAddressPostalCodeMoreThan10Chars
	}

	return nil
}
