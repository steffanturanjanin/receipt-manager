package models

import (
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
)

type User struct {
	ID        uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	FirstName string    `gorm:"not null; size:255" json:"firstName"`
	LastName  string    `gorm:"not null; size:255" json:"lastName"`
	Email     string    `gorm:"unique; not null; size:100" json:"email"`
	Password  string    `gorm:"not null; size:100"`
	CreatedAt time.Time `gorm:"not null; autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null; autoCreateTime" json:"updatedAt"`
}

func NewUserResponseDTOFromUserModel(user User) dto.User {
	return dto.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserModelFromRegisterRequestDTO(requestDTO dto.RegisterUserRequest) User {
	return User{
		FirstName: requestDTO.FirstName,
		LastName:  requestDTO.LastName,
		Email:     requestDTO.Email,
		Password:  requestDTO.Password,
	}
}
