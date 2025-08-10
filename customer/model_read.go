package customer

import (
	"time"

	"github.com/mmiftahrzki/customer/customer/address"
)

type CustomerReadModel struct {
	Id        uint8           `json:"id"`
	Email     string          `json:"email" validate:"required, email,max=100"`
	FullName  string          `json:"full_name" validate:"required,max=255"`
	Address   address.Address `json:"address"`
	CreatedAt time.Time       `json:"created_at"`
}
