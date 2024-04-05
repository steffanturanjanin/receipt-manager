package database

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

var (
	categories = []struct {
		Name  string
		Color string
	}{
		{Name: "Voće i povrće", Color: "#ec3f7e"},
		{Name: "Mesne prerađevine", Color: "#578f5b"},
		{Name: "Meso i riba", Color: "#a129d7"},
		{Name: "Mlečni proizvodi i jaja", Color: "#c8fedd"},
		{Name: "Hleb i pekarski proizvodi", Color: "#5a75db"},
		{Name: "Namirnice za pripremu jela", Color: "#0f06b0"},
		{Name: "Slatkiši i grickalice", Color: "#cac46d"},
		{Name: "Gotova jela", Color: "#9613c5"},
		{Name: "Alkoholna pića", Color: "#2dc668"},
		{Name: "Sokovi i voda", Color: "#eaf991"},
		{Name: "Čaj i kafa", Color: "#4f3531"},
		{Name: "Kućna hemija", Color: "#057b41"},
		{Name: "Lična higijena i nega", Color: "#ecdb08"},
		{Name: "Tehnika", Color: "#850312"},
		{Name: "Odeća i obuća", Color: "#e7e6db"},
		{Name: "Zdravlje", Color: "#f94c79"},
		{Name: "Domaćinstvo", Color: "#9becec"},
		{Name: "Kućni ljubimci", Color: "#194e36"},
		{Name: "Automobil", Color: "#e11342"},
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
		if err := createCategory(db, category.Name, category.Color); err != nil {
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

func createCategory(db *gorm.DB, name string, color string) error {
	return db.Create(&models.Category{Name: name, Color: color}).Error
}

func createTax(db *gorm.DB, tax *models.Tax) error {
	return db.Create(&tax).Error
}
