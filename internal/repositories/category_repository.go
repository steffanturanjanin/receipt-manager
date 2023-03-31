package repositories

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	GetAll() ([]models.Category, error)
	GetIds() ([]int, error)
	GetById(int) (*models.Category, error)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetIds() ([]int, error) {
	ids := make([]int, 0)

	if err := r.db.Model(&models.Category{}).Select("id").Find(&ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *CategoryRepository) GetById(id int) (*models.Category, error) {
	var category models.Category

	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
