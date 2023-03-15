package models

import "time"

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
