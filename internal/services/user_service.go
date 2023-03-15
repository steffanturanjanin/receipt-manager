package services

import (
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
)

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserService(repository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: repository,
	}
}

func (service *UserService) Create(request dto.RegisterUserRequest) (*dto.User, error) {
	userModel, err := service.UserRepository.Create(request)
	if err != nil {
		return nil, err
	}

	response := models.NewUserResponseDTOFromUserModel(*userModel)

	return &response, nil
}

func (service *UserService) GetById(id int) (*dto.User, error) {
	userModel, err := service.UserRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	userResponse := models.NewUserResponseDTOFromUserModel(*userModel)

	return &userResponse, nil
}
