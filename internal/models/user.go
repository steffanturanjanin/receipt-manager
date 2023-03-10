package models

import (
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

type User struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	FirstName string    `gorm:"size:255;not null" json:"first_name"`
	LastName  string    `gorm:"size:255;not null" json:"last_name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoCreateTime" json:"updated_at"`
}

func NewUserResponseDTOFromUserModel(user User) transport.UserResponseDTO {
	return transport.UserResponseDTO{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserModelFromRegisterRequestDTO(requestDTO transport.RegisterUserRequestDTO) User {
	return User{
		FirstName: requestDTO.FirstName,
		LastName:  requestDTO.LastName,
		Email:     requestDTO.Email,
		Password:  requestDTO.Password,
	}
}
