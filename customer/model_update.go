package customer

import (
	"fmt"
)

type CustomerUpdateModel struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
}

var errCustomerFirstNameNull = fmt.Errorf("customer first name is required")
var errCustomerLastNameNull = fmt.Errorf("customer last name is required")

func ValidateCustomerUpdateModel(m CustomerUpdateModel) error {
	if m.FirstName == nil || len(*m.FirstName) > 45 {
		return errCustomerFirstNameNull
	}

	if m.LastName == nil || len(*m.LastName) > 45 {
		return errCustomerLastNameNull
	}

	return nil
}
