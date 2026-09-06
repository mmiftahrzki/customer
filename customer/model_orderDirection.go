package customer

import (
	"errors"
	"strings"
)

type orderDirection string

const (
	Ascending  orderDirection = "ASC"
	Descending orderDirection = "DESC"
)

func ParseOrderDirection(s string) (orderDirection, error) {
	clean := strings.ToLower(strings.TrimSpace(s))

	switch clean {
	case "asc", "":
		return Ascending, nil
	case "desc":
		return Descending, nil
	default:
		return "", errors.New("invalid direction: must be 'ASC' or 'DESC'")
	}
}
