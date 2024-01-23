package models

import (
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
)

type Store struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Tin          string    `gorm:"not null; type:varchar(9)" json:"tin"`
	Name         string    `gorm:"not null; size:255" json:"name"`
	LocationId   string    `gorm:"not null; size:255" json:"locationId"`
	LocationName string    `gorm:"not null; size:255" json:"locationName"`
	Address      string    `gorm:"not null; size:255" json:"address"`
	City         string    `gorm:"not null; size:255" json:"city"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoCreateTime" json:"updatedAt"`
	Receipts     []Receipt `gorm:"foreignKey:StoreID; references:ID; constraint:OnDelete:SET NULL" json:"receipts"`
}

func NewStoreFromStoreDTO(storeDTO dto.Store) Store {
	return Store{
		Tin:          storeDTO.Tin,
		Name:         storeDTO.Name,
		LocationId:   storeDTO.LocationId,
		LocationName: storeDTO.Name,
		Address:      storeDTO.Address,
		City:         storeDTO.City,
	}
}

func (store Store) NewStoreDTO() dto.Store {
	return dto.Store{
		Tin:          store.Tin,
		Name:         store.Name,
		LocationName: store.LocationName,
		LocationId:   store.LocationId,
		Address:      store.Address,
		City:         store.City,
	}
}
