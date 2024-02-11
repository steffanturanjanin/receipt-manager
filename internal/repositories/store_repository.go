package repositories

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/pkg/receipt-fetcher"
	"gorm.io/gorm"
)

type StoreRepositoryInterface interface {
	GetByTin(string) (*models.Store, error)
	FirstOrCreateFromDto(storeDto receipt_fetcher.Store) (*models.Store, error)
}

type StoreRepository struct {
	DB *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{
		DB: db,
	}
}

func (r *StoreRepository) GetByTin(tin string) (*models.Store, error) {
	var store *models.Store
	if err := r.DB.First(&store, tin).Error; err != nil {
		return nil, err
	}

	return store, nil
}

func (r *StoreRepository) FirstOrCreateFromDto(storeDto receipt_fetcher.Store) (*models.Store, error) {
	store := models.Store{
		Tin:          storeDto.Tin,
		Name:         storeDto.Name,
		LocationName: storeDto.LocationName,
		LocationId:   storeDto.LocationId,
		Address:      storeDto.Address,
		City:         storeDto.City,
	}

	r.DB.Where(models.Store{Tin: storeDto.Tin}).FirstOrCreate(&store)

	return &store, nil
}
