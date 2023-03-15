package database

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Receipt{},
		&models.ReceiptItem{},
		&models.Tax{},
		&models.Store{},
	)
}
