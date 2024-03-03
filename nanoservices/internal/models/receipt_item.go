package models

import (
	"math"
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
)

type ReceiptItem struct {
	ID           uint      `gorm:"primaryKey; autoIncrement" json:"id"`
	ReceiptID    uint      `gorm:"not null" json:"receiptId"`
	CategoryID   *uint     `gorm:"nullable" json:"categoryId"`
	TaxID        *uint     `gorm:"nullable" json:"taxId"`
	Name         string    `gorm:"not null; size:255" json:"name"`
	Unit         string    `gorm:"not null" json:"unit"`
	Quantity     float64   `gorm:"not null" json:"quantity"`
	SingleAmount int       `gorm:"not null; default:0" json:"singleAmount"`
	TotalAmount  int       `gorm:"not null; default:0" json:"totalAmount"`
	CreatedAt    time.Time `gorm:"not null; autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"not null; autoCreateTime" json:"updatedAt"`
	Receipt      Receipt   `gorm:"foreignKey:ReceiptID; references:ID" json:"receipt"`
	Category     *Category `gorm:"foreignKey:CategoryID; references:ID" json:"category"`
	Tax          *Tax      `gorm:"foreignKey:TaxID; references:ID" json:"tax"`
}

func (receiptItem ReceiptItem) NewReceiptItemDTO() dto.ReceiptItem {
	return dto.ReceiptItem{
		ID:       receiptItem.ID,
		Name:     receiptItem.Name,
		Quantity: receiptItem.Quantity,
		Unit:     receiptItem.Unit,
		//Tax:          *dto.TaxIdentifier(receiptItem.Tax).Tax(),
		SingleAmount: math.Round(float64(receiptItem.SingleAmount)) / 100,
		TotalAmount:  math.Round(float64(receiptItem.TotalAmount)) / 100,
	}
}
