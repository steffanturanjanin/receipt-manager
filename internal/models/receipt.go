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
	ID                  uint           `gorm:"primaryKey; autoIncrement" json:"id"`
	Status              string         `gorm:"not null" json:"status"`
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

func (r Receipt) NewReceiptDTO() (*dto.Receipt, error) {
	receiptItems := make([]dto.ReceiptItem, 0)
	for _, receiptItemModel := range r.ReceiptItems {
		receiptItem := receiptItemModel.NewReceiptItemDTO()
		receiptItems = append(receiptItems, receiptItem)
	}

	taxes := make([]dto.Tax, 0)
	for _, taxModel := range r.Taxes {
		if tax := taxModel.NewTaxDTO(); tax != nil {
			taxes = append(taxes, *tax)
		}
	}

	meta := make(map[string]string)
	if err := json.Unmarshal(r.Meta, &meta); err != nil {
		return nil, err
	}

	receiptDTO := dto.Receipt{
		ID:                  r.ID,
		Store:               r.Store.NewStoreDTO(),
		PfrNumber:           r.PfrNumber,
		Counter:             r.Counter,
		TotalPurchaseAmount: math.Round(float64(r.TotalPurchaseAmount)) / 100,
		TotalTaxAmount:      math.Round(float64(r.TotalTaxAmount)) / 100,
		ReceiptItems:        receiptItems,
		Taxes:               taxes,
		Date:                r.Date,
		QrCode:              r.QrCode,
		Meta:                meta,
		CreatedAt:           r.Date,
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
