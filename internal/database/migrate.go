package database

import (
	"errors"

	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Category{}); err == nil && db.Migrator().HasTable(&models.Category{}) {
		if err := db.First(&models.Category{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Run categories seeder
			SeedCategories(db)
		}
	}

	if err := db.AutoMigrate(&models.Tax{}); err == nil && db.Migrator().HasTable(&models.Tax{}) {
		if err := db.First(&models.Tax{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Run taxes seeder
			SeedTaxes(db)
		}
	}

	return db.AutoMigrate(
		&models.User{},
		&models.Store{},
		&models.Receipt{},
		&models.ReceiptItem{},
	)
}
