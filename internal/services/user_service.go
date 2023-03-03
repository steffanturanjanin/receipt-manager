package services

import (
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/transform"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
)

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserService(repository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: repository,
	}
}

func (service *UserService) Create(request transport.RegisterUserRequest) (*transport.UserResponse, error) {
	userModel, err := service.UserRepository.Create(request)
	if err != nil {
		return nil, err
	}

	response := transform.NewUserResponseFromUserModel(*userModel)

	return &response, nil
}

func (service *UserService) GetById(id int) (*transport.UserResponse, error) {
	userModel, err := service.UserRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	userResponse := transform.NewUserResponseFromUserModel(*userModel)

	return &userResponse, nil
}
