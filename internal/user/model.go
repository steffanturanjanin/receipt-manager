package user

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string    `gorm:"size:255;not null" json:"first_name"`
	LastName  string    `gorm:"size:255;not null" json:"last_name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoCreateTime" json:"updated_at"`
}

func NewUserResponseDTOFromUserModel(user User) UserDto {
	return UserDto{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (user *User) FromRegisterRequest(registerRequest RegisterUserRequest) *User {
	user.FirstName = registerRequest.FirstName
	user.LastName = registerRequest.LastName
	user.Email = registerRequest.Email
	user.Password = registerRequest.Password

	return user
}

func NewUserModelFromRegisterRequestDTO(requestDTO RegisterUserRequest) User {
	return User{
		FirstName: requestDTO.FirstName,
		LastName:  requestDTO.LastName,
		Email:     requestDTO.Email,
		Password:  requestDTO.Password,
	}
}
