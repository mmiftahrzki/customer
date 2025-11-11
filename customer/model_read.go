package customer

import (
	"strings"
	"time"

	"github.com/mmiftahrzki/customer/customer/address"
)

type customerAddressReadModel string
type customerReadModel struct {
	Id        int                      `json:"id"`
	Email     string                   `json:"email" validate:"required, email,max=100"`
	FullName  string                   `json:"full_name" validate:"required,max=255"`
	Address   customerAddressReadModel `json:"address"`
	CreatedAt time.Time                `json:"created_at"`
}

func newCustomerReadModel(sql_model customerSQLModel) (customer customerReadModel) {
	address_read_model := address.AddressReadModel{}

	if sql_model.customer_id.Valid {
		customer.Id = int(sql_model.customer_id.Int16)
	}

	if sql_model.email.Valid {
		customer.Email = sql_model.email.String
	}

	if sql_model.first_name.Valid {
		customer.FullName = sql_model.first_name.String
	}

	if sql_model.last_name.Valid && sql_model.last_name.String != "" {
		customer.FullName += " " + sql_model.last_name.String
	}

	if sql_model.address_id.Valid {
		address_read_model.Id = int(sql_model.address_id.Int16)
	}

	if sql_model.address.Address.Valid {
		address_read_model.Address = sql_model.address.Address.String
	}

	if sql_model.address.Address2.Valid {
		address_read_model.Address2 = sql_model.address.Address2.String
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

	if sql_model.create_date.Valid {
		customer.CreatedAt = sql_model.create_date.Time
	}

	return
}

func newAddressReadModel(address address.AddressReadModel) customerAddressReadModel {
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

	return customerAddressReadModel(address_str)
}
