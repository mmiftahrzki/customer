package customer

import (
	"strings"
	"time"

	"github.com/mmiftahrzki/customer/customer/address"
)

type modelReadAddress string
type ModelRead struct {
	Id        int              `json:"id"`
	Email     string           `json:"email" validate:"required, email,max=100"`
	FullName  string           `json:"full_name" validate:"required,max=255"`
	Address   modelReadAddress `json:"address"`
	CreatedAt time.Time        `json:"created_at"`
}

func newReadModel(modelSQL modelSQL) ModelRead {
	var customer ModelRead
	var addressModelRead address.ModelRead

	if modelSQL.id.Valid {
		customer.Id = int(modelSQL.id.Int16)
	}

	if modelSQL.email.Valid {
		customer.Email = modelSQL.email.String
	}

	if modelSQL.firstName.Valid {
		customer.FullName = modelSQL.firstName.String
	}

	if modelSQL.lastName.Valid && modelSQL.lastName.String != "" {
		customer.FullName += " " + modelSQL.lastName.String
	}

	if modelSQL.addressId.Valid {
		addressModelRead.Id = int(modelSQL.addressId.Int16)
	}

	if modelSQL.address.Address.Valid {
		addressModelRead.Address = modelSQL.address.Address.String
	}

	if modelSQL.address.District.Valid {
		addressModelRead.District = modelSQL.address.District.String
	}

	if modelSQL.address.CityId.Valid {
		addressModelRead.CityId = int(modelSQL.address.CityId.Int16)
		addressModelRead.City = modelSQL.address.City.String
	}

	if modelSQL.address.PostalCode.Valid {
		addressModelRead.PostalCode = modelSQL.address.PostalCode.String
	}

	if modelSQL.address.CountryId.Valid {
		addressModelRead.CountryId = int(modelSQL.address.CountryId.Int16)
		addressModelRead.Country = modelSQL.address.Country.String
	}

	customer.Address = newReadModelAddress(addressModelRead)

	if modelSQL.createdAt.Valid {
		customer.CreatedAt = modelSQL.createdAt.Time
	}

	return customer
}

func newReadModelAddress(address address.ModelRead) modelReadAddress {
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

	if address.CityId != 0 {
		addresses = append(addresses, address.City)
	}

	if address.CountryId != 0 {
		addresses = append(addresses, address.Country)
	}

	addressStr := strings.Join(addresses, " ")

	return modelReadAddress(addressStr)
}
