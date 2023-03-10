package dto

import (
	"time"
)

type RegisterUserRequestDTO struct {
	FirstName string `validate:"required,max=255" json:"first_name"`
	LastName  string `validate:"required,max=255" json:"last_name"`
	Email     string `validate:"required,email,max=100" json:"email"`
	Password  string `validate:"required,min=8,max=100" json:"password"`
}

type LoginUserRequestDTO struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

type UserResponseDTO struct {
	Id        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginUserResponseDTO struct {
	AccessToken string `json:"access_token"`
}
