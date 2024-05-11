package models

import (
	"encoding/json"
	"math"
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"gorm.io/datatypes"
)

const (
	RECEIPT_STATUS_PENDING   = "pending"
	RECEIPT_STATUS_PROCESSED = "processed"
	RECEIPT_STATUS_FAILED    = "failed"
)

type Receipt struct {
	ID                  uint            `gorm:"primaryKey; autoIncrement" json:"id"`
	Vl                  string          `gorm:"type:text; not null" json:"vl"`
	UserID              *uint           `gorm:"nullable" json:"userId"`
	StoreID             *uint           `gorm:"nullable" json:"storeId"`
	Status              string          `gorm:"not null" json:"status"`
	PfrNumber           *string         `gorm:"unique; nullable" json:"pfrNumber"`
	Counter             *string         `gorm:"unique; nullable" json:"counter"`
	TotalPurchaseAmount *int            `gorm:"nullable; default:0" json:"totalPurchaseAmount"`
	TotalTaxAmount      *int            `gorm:"nullable; default:0" json:"totalTaxAmount"`
	Date                *time.Time      `gorm:"nullable" json:"date"`
	QrCode              *string         `gorm:"type:text; nullable" json:"qrCode"`
	Meta                *datatypes.JSON `gorm:"nullable" json:"metaData"`
	IsFavorite          bool            `gorm:"default:false" json:"isFavorite"`
	CreatedAt           time.Time       `gorm:"not null; autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time       `gorm:"not null; autoCreateTime" json:"updatedAt"`
	User                *User           `gorm:"foreignKey:UserID; references:ID; nullable" json:"user"`
	Store               *Store          `gorm:"foreignKey:StoreID; references:ID; nullable" json:"store"`
	ReceiptItems        []ReceiptItem   `gorm:"foreignKey:ReceiptID; references:ID; constraint:OnDelete:CASCADE; nullable" json:"receiptItems"`
}

func (r Receipt) NewReceiptDTO() (*dto.Receipt, error) {
	receiptItems := make([]dto.ReceiptItem, 0)
	for _, receiptItemModel := range r.ReceiptItems {
		receiptItem := receiptItemModel.NewReceiptItemDTO()
		receiptItems = append(receiptItems, receiptItem)
	}

	// taxes := make([]dto.Tax, 0)
	// for _, taxModel := range r.Taxes {
	// 	if tax := taxModel.NewTaxDTO(); tax != nil {
	// 		taxes = append(taxes, *tax)
	// 	}
	// }

	meta := make(map[string]string)
	if err := json.Unmarshal(*r.Meta, &meta); err != nil {
		return nil, err
	}

	receiptDTO := dto.Receipt{
		ID:                  r.ID,
		Store:               r.Store.NewStoreDTO(),
		PfrNumber:           *r.PfrNumber,
		Counter:             *r.Counter,
		TotalPurchaseAmount: math.Round(float64(*r.TotalPurchaseAmount)) / 100,
		TotalTaxAmount:      math.Round(float64(*r.TotalTaxAmount)) / 100,
		ReceiptItems:        receiptItems,
		Date:                *r.Date,
		QrCode:              *r.QrCode,
		Meta:                meta,
		CreatedAt:           *r.Date,
	}

	return &receiptDTO, nil
}

func (r Receipt) SortableFields() []string {
	return []string{
		"id",
		"pfr_number",
		"counter",
		"total_purchase_amount",
		"total_tax_amount",
		"date",
	}
}
