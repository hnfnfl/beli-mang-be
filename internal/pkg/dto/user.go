package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

type Sort string

const (
	ASC  Sort = "ASC"
	DESC Sort = "DESC"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.Length(5, 30)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(5, 30)),
	)
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.Length(5, 30)),
		validation.Field(&r.Password, validation.Required, validation.Length(5, 30)),
	)
}
