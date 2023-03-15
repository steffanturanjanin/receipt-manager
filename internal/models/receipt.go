package models

import (
	"time"

	"gorm.io/datatypes"
)

type Receipt struct {
	ID                  uint           `gorm:"primaryKey; autoIncrement" json:"id"`
	PfrNumber           string         `gorm:"unique; not null" json:"pfr_number"`
	Counter             string         `gorm:"unique; not null" json:"counter"`
	TotalPurchaseAmount int            `gorm:"not null; default:0" json:"total_purchase_amount"`
	TotalTaxAmount      int            `gorm:"not null; default:0" json:"total_tax_amount"`
	Date                time.Time      `gorm:"not null" json:"date"`
	QrCode              string         `gorm:"not null" json:"qr_code"`
	CreatedAt           time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	StoreID             string         `gorm:"type:varchar(9)"`
	ReceiptItems        []ReceiptItem  `gorm:"foreignKey:ReceiptID;references:ID;constraint:OnDelete:CASCADE" json:"receipt_items"`
	Taxes               []Tax          `gorm:"foreignKey:ReceiptID;references:ID;constraint:OnDelete:CASCADE" json:"taxes"`
	Meta                datatypes.JSON `gorm:"nullable" json:"meta_data"`
	Store               Store          `json:"store"`
}
