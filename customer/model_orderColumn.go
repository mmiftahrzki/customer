package customer

import (
	"errors"
	"strings"
)

type orderColumn string

const (
	Name    orderColumn = "first_name, last_name"
	Email   orderColumn = "email"
	Address             = "address"
)

func ParseOrderColumn(s string) (orderColumn, error) {
	clean := strings.ToLower(strings.TrimSpace(s))

	switch clean {
	case "name", "":
		return Name, nil
	case "email":
		return Email, nil
	case "address":
		return Address, nil
	default:
		return "", errors.New("invalid order column: must be 'name', 'email' or 'address'")
	}
}
