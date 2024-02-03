package user

import (
	"time"
)

type User struct {
	ID             string    `json:"id" readOnly:"true"`
	FirstName      string    `json:"first_name" validate:"required"`
	LastName       string    `json:"last_name" validate:"required"`
	Password       string    `json:"password" validate:"required" writeOnly:"true"`
	Username       string    `json:"username" validate:"required,email"`
	AccountCreated time.Time `json:"account_created" readOnly:"true"`
	AccountUpdated time.Time `json:"account_updated" readOnly:"true"`
}
