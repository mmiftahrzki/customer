package customer

import (
	"errors"
)

type customerUpdateModel struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
}

var errCustomerFirstNameNull = errors.New("customer first name is required")
var errCustomerLastNameNull = errors.New("customer last name is required")

func validatecustomerUpdateModel(m customerUpdateModel) error {
	if m.FirstName == nil || len(*m.FirstName) > 45 {
		return errCustomerFirstNameNull
	}

	if m.LastName == nil || len(*m.LastName) > 45 {
		return errCustomerLastNameNull
	}

	return nil
}
