package models

import "github.com/go-playground/validator/v10"

var (
	CustomerRole = "customer"
	WorkerRole   = "worker"

	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

type RegisterUserRequest struct {
	Username string `db:"username" json:"username" validate:"required,gte=2"`
	Password string `db:"password" json:"password" validate:"required,gte=6"`
	Role     string `db:"role" json:"role" validate:"required"`
}

func (i *RegisterUserRequest) Validate() error {
	return validate.Struct(i)
}

type LoginUserRequest struct {
	Username string `db:"username" json:"username" validate:"required,gte=2"`
	Password string `db:"password" json:"password" validate:"required,gte=6"`
}

func (i *LoginUserRequest) Validate() error {
	return validate.Struct(i)
}
