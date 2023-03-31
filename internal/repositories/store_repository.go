package repositories

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

type StoreRepositoryInterface interface {
	GetByTin(string) (*models.Store, error)
}

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

func (r *StoreRepository) GetByTin(tin string) (*models.Store, error) {
	var store *models.Store
	if err := r.db.First(&store, tin).Error; err != nil {
		return nil, err
	}

	return store, nil
}
