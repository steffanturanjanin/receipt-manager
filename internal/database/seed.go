package database

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

var (
	categories = []string{
		"Hrana",
		"Voće i povrće",
		"Delikates",
		"Meso",
		"Riba",
		"Mlečni proizvodi i jaja",
		"Hleb i pekarski proizvodi",
		"Namirnice za pripremu jela",
		"Slatkiši i grickalice",
		"Gotova jela",
		"Piće",
		"Alkoholna pića",
		"Sokovi",
		"Voda",
		"Čaj i kafa",
		"Kućna hemija",
		"Lična higijena i nega",
		"Tehnika",
		"Odeća i obuća",
		"Zdravlje",
		"Domaćinstvo",
	}

	taxes = []models.Tax{
		{
			Identifier: "Ђ",
			Name:       "О-ПДВ",
			Rate:       20,
		},
		{
			Identifier: "Е",
			Name:       "П-ПДВ",
			Rate:       10,
		},
	}
)

func SeedCategories(db *gorm.DB) error {
	for _, category := range categories {
		if err := createCategory(db, category); err != nil {
			return err
		}
	}

	return nil
}

func SeedTaxes(db *gorm.DB) error {
	for _, tax := range taxes {
		if err := createTax(db, &tax); err != nil {
			return err
		}
	}

	return nil
}

func createCategory(db *gorm.DB, name string) error {
	return db.Create(&models.Category{Name: name}).Error
}

func createTax(db *gorm.DB, tax *models.Tax) error {
	return db.Create(&tax).Error
}
