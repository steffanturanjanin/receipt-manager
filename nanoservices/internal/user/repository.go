package user

import (
	native_errors "errors"

	"github.com/go-sql-driver/mysql"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(request RegisterUserRequest) (*User, error)
	GetByEmail(email string) (*User, error)
	GetById(id int) (*User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) Create(request RegisterUserRequest) (*User, error) {
	userModel := new(User)
	userModel.FromRegisterRequest(request)

	hashedPassword, err := utils.HashPassword(userModel.Password)
	if err != nil {
		return nil, err
	}

	userModel.Password = hashedPassword
	result := repository.db.Create(&userModel)

	if result.Error != nil {
		err := result.Error
		var mysqlErr *mysql.MySQLError
		if native_errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			err = errors.NewErrDuplicateEntry(err, "User with given email already exists.")
		}

		return nil, err
	}

	return userModel, nil
}

func (repository *UserRepository) GetByEmail(email string) (*User, error) {
	var userModel *User
	result := repository.db.First(&userModel, "email = ?", email)

	if result.Error != nil {
		if native_errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewErrResourceNotFound(result.Error, "User not found.")
		}
	}

	return userModel, nil
}

func (repository *UserRepository) GetById(id int) (*User, error) {
	var userModel *User
	result := repository.db.First(&userModel, "id = ?", id)

	if result.Error != nil {
		if native_errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NewErrResourceNotFound(result.Error, "User not found.")
		}
	}

	return userModel, nil
}
