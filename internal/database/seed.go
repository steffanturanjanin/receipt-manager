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
)

func SeedCategories(db *gorm.DB) error {
	for _, category := range categories {
		if err := createCategory(db, category); err != nil {
			return err
		}
	}

	return nil
}

// type Seed struct {
// 	Run func(*gorm.DB) error
// }

// type SeederInterface interface {
// 	Seed(*gorm.DB)
// }

// type seeder struct {
// 	seeds []Seed
// }

// func (s seeder) Seed(db *gorm.DB) {
// 	for _, seed := range s.seeds {
// 		seed.Run(db)
// 	}
// }

// type CategorySeeder struct {
// 	seeder
// }

// func NewCategorySeeder() CategorySeeder {
// 	categorySeeds := make([]Seed, 0)

// 	for _, category := range categories {
// 		categorySeeds = append(categorySeeds, Seed{
// 			Run: func(db *gorm.DB) error {
// 				return createCategory(db, category)
// 			},
// 		})
// 	}

// 	return CategorySeeder{seeder{seeds: categorySeeds}}
// }

func createCategory(db *gorm.DB, name string) error {
	return db.Create(&models.Category{Name: name}).Error
}
