package repositories

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

type ReceiptItemRepositoryInterface interface {
	Update(*models.ReceiptItem) error
	FindById(int) (*models.ReceiptItem, error)
}

type ReceiptItemRepository struct {
	db *gorm.DB
}

func NewReceiptItemRepository(db *gorm.DB) *ReceiptItemRepository {
	return &ReceiptItemRepository{
		db: db,
	}
}

func (r *ReceiptItemRepository) Update(receiptItem *models.ReceiptItem) error {
	return r.db.Updates(&receiptItem).Error
}

func (r *ReceiptItemRepository) FindById(id int) (*models.ReceiptItem, error) {
	var receiptItem *models.ReceiptItem

	if err := r.db.Preload("Category").First(&receiptItem, uint(id)).Error; err != nil {
		return nil, err
	}

	return receiptItem, nil
}
