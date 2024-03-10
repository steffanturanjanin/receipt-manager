package dto

type RegisterUserRequest struct {
	FirstName string `validate:"required,max=255" json:"first_name"`
	LastName  string `validate:"required,max=255" json:"last_name"`
	Email     string `validate:"required,email,max=100" json:"email"`
	Password  string `validate:"required,min=8,max=100" json:"password"`
}

type LoginUserRequest struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

type User struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}
