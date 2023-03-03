package repositories

import (
	"errors"

	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/transform"
	"github.com/steffanturanjanin/receipt-manager/internal/transport"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(request transport.RegisterUserRequest) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetById(id int) (*models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) Create(request transport.RegisterUserRequest) (*models.User, error) {
	userModel := transform.NewUserModelFromRegisterRequest(request)
	hashedPassword, err := utils.HashPassword(userModel.Password)
	if err != nil {
		return nil, err
	}

	userModel.Password = hashedPassword
	result := repository.db.Create(&userModel)

	if result.Error != nil {
		return nil, result.Error
	}

	return &userModel, nil
}

func (repository *UserRepository) GetByEmail(email string) (*models.User, error) {
	var userModel *models.User
	repository.db.First(&userModel, "email = ?", email)

	if userModel == nil {
		return nil, errors.New("user with requested email not found")
	}

	return userModel, nil
}

func (repository *UserRepository) GetById(id int) (*models.User, error) {
	var userModel *models.User
	repository.db.First(&userModel, "id = ?", id)

	if userModel == nil {
		return nil, errors.New("user with requested id not found")
	}

	return userModel, nil
}
