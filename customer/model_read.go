package customer

import (
	"strings"
	"time"

	"github.com/mmiftahrzki/customer/customer/address"
)

type addressReadModel string
type modelRead struct {
	Id        int              `json:"id"`
	Email     string           `json:"email" validate:"required, email,max=100"`
	FullName  string           `json:"full_name" validate:"required,max=255"`
	Address   addressReadModel `json:"address"`
	CreatedAt time.Time        `json:"created_at"`
}

func newReadModel(sql_model modelSQL) (customer modelRead) {
	address_read_model := address.ModelRead{}

	if sql_model.id.Valid {
		customer.Id = int(sql_model.id.Int16)
	}

	if sql_model.email.Valid {
		customer.Email = sql_model.email.String
	}

	if sql_model.firstName.Valid {
		customer.FullName = sql_model.firstName.String
	}

	if sql_model.lastName.Valid && sql_model.lastName.String != "" {
		customer.FullName += " " + sql_model.lastName.String
	}

	if sql_model.addressId.Valid {
		address_read_model.Id = int(sql_model.addressId.Int16)
	}

	if sql_model.address.Address.Valid {
		address_read_model.Address = sql_model.address.Address.String
	}

	if sql_model.address.District.Valid {
		address_read_model.District = sql_model.address.District.String
	}

	if sql_model.address.CityId.Valid {
		address_read_model.CityId = int(sql_model.address.CityId.Int16)
	}

	if sql_model.address.PostalCode.Valid {
		address_read_model.PostalCode = sql_model.address.PostalCode.String
	}

	customer.Address = newAddressReadModel(address_read_model)

	if sql_model.createdAt.Valid {
		customer.CreatedAt = sql_model.createdAt.Time
	}

	return
}

func newAddressReadModel(address address.ModelRead) addressReadModel {
	addresses := []string{}

	if address.Address != "" {
		addresses = append(addresses, address.Address)
	}

	if address.Address2 != "" {
		addresses = append(addresses, address.Address2)
	}

	if address.District != "" {
		addresses = append(addresses, address.District)
	}

	if address.PostalCode != "" {
		addresses = append(addresses, address.PostalCode)
	}

	address_str := strings.Join(addresses, " ")

	return addressReadModel(address_str)
}
