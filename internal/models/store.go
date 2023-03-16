package models

import (
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
)

type Store struct {
	Tin          string    `gorm:"type:varchar(9);primaryKey" json:"tin"`
	Name         string    `gorm:"not null;size:255" json:"name"`
	LocationId   string    `gorm:"not null;size:255" json:"location_id"`
	LocationName string    `gorm:"not null;size:255" json:"location_name"`
	Address      string    `gorm:"not null;size:255" json:"address"`
	City         string    `gorm:"not null;size:255" json:"city"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	Receipts     []Receipt `gorm:"foreignKey:StoreID;references:Tin;constraint:OnDelete:SET NULL" json:"receipts"`
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
