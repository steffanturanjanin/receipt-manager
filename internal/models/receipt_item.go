package models

import (
	"math"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
)

type ReceiptItem struct {
	ID           uint    `gorm:"primaryKey; autoIncrement" json:"id"`
	Name         string  `gorm:"size:255,not null" json:"name"`
	Unit         string  `gorm:"not null" json:"unit"`
	Quantity     float64 `gorm:"not null" json:"quantity"`
	SingleAmount int     `gorm:"not null; default:0" json:"single_amount"`
	TotalAmount  int     `gorm:"not null; default:0" json:"total_amount"`
	Tax          int     `gorm:"nullable" json:"tax"`
	ReceiptID    uint    `json:"receipt_id"`
	CategoryID   uint    `json:"category_id"`
}

func (receiptItem ReceiptItem) NewReceiptItemDTO() dto.ReceiptItem {
	return dto.ReceiptItem{
		ID:           receiptItem.ID,
		Name:         receiptItem.Name,
		Quantity:     receiptItem.Quantity,
		Unit:         receiptItem.Unit,
		Tax:          *dto.TaxIdentifier(receiptItem.Tax).Tax(),
		SingleAmount: math.Round(float64(receiptItem.SingleAmount)) / 100,
		TotalAmount:  math.Round(float64(receiptItem.TotalAmount)) / 100,
	}
}
