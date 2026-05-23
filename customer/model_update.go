package customer

import (
	"encoding/json"
	"errors"

	"github.com/mmiftahrzki/customer/customer/address"
)

type modelUpdate struct {
	FirstName *string              `json:"first_name"`
	LastName  *string              `json:"last_name"`
	Email     *string              `json:"email"`
	Address   *address.ModelUpdate `json:"address"`
}

func (this *modelUpdate) UnmarshalJSON(data []byte) error {
	type Alias modelUpdate

	var model Alias

	err := json.Unmarshal(data, &model)
	if err != nil {
		return err
	}

	if model.FirstName == nil {
		return errors.New("first name is nil")
	}

	if model.LastName == nil {
		return errors.New("last name is nil")
	}

	if model.Email == nil {
		return errors.New("email is nil")
	}

	*this = modelUpdate(model)

	return nil
}

func (this *modelUpdate) Validate() error {
	if len(*this.FirstName) > 45 {
		return errors.New("first name cannot be more than 45 characters")
	}

	if len(*this.LastName) > 45 {
		return errors.New("last name cannot be more than 45 characters")
	}

	if len(*this.Email) > 50 {
		return errors.New("last name cannot be more than 50 characters")
	}

	return nil
}
