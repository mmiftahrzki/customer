package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email" validate:"required,email,max=100"`
	Password  string    `json:"password,omitempty" validate:"required,max=32"`
	Fullname  string    `json:"fullname" validate:"required,max=255"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by,omitempty"`
}
