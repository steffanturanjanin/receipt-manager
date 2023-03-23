package database

import (
	"errors"

	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Category{}); err == nil && db.Migrator().HasTable(&models.Category{}) {
		if err := db.First(&models.Category{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert seed data
			SeedCategories(db)
		}
	}

	return db.AutoMigrate(
		&models.User{},
		&models.Receipt{},
		&models.ReceiptItem{},
		&models.Tax{},
		&models.Store{},
	)
}
