package transport

import (
	"time"
)

type RegisterUserRequest struct {
	FirstName string `validate:"required" json:"first_name"`
	LastName  string `validate:"required" json:"last_name"`
	Email     string `validate:"required,email" json:"email"`
	Password  string `validate:"required" json:"password"`
}

type LoginUserRequest struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

type UserResponse struct {
	Id        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}
