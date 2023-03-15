package dto

import (
	"time"

	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

type ReceiptData receipt_fetcher.Receipt

// ReceiptDTO
type Receipt struct {
	ID                  uint              `json:"id"`
	Store               Store             `json:"store"`
	PfrNumber           string            `json:"pfr_number"`
	Counter             string            `json:"counter"`
	TotalPurchaseAmount float64           `json:"total_purchase_amount"`
	TotalTaxAmount      float64           `json:"total_tax_amount"`
	ReceiptItems        []ReceiptItem     `json:"receipt_items"`
	Taxes               []Tax             `json:"taxes"`
	Date                time.Time         `json:"date"`
	QrCode              string            `json:"qr_code"`
	Meta                map[string]string `json:"meta_data"`
	CreatedAt           time.Time         `json:"created_at"`
}

// StoreDTO
type Store struct {
	Tin          string `jsontin:"tin"`
	Name         string `json:"name"`
	LocationName string `json:"location_name"`
	LocationId   string `json:"location_id"`
	Address      string `json:"address"`
	City         string `json:"city"`
}

// ReceiptItem DTO
type ReceiptItem struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Tax          Tax     `json:"tax"`
	SingleAmount float64 `json:"single_amount"`
	TotalAmount  float64 `json:"total_amount"`
}

type Tax struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Rate       int    `json:"rate"`
}

type TaxIdentifier int

const (
	TaxPPdvNumericalIdentifier TaxIdentifier = iota
	TaxOPdvNumericalIdentifier

	TaxPPdvAlphaIdentifier string = "Ђ"
	TaxOPdvAlphaIdentifier string = "Е"
)

var (
	TaxPPdv = Tax{Name: "О-ПДВ", Identifier: TaxPPdvAlphaIdentifier, Rate: 20}
	TaxOPdv = Tax{Name: "П-ПДВ", Identifier: TaxOPdvAlphaIdentifier, Rate: 10}

	AllowedTaxNumericalIdentifierValues = []TaxIdentifier{
		TaxPPdvNumericalIdentifier,
		TaxOPdvNumericalIdentifier,
	}

	TaxIdentifierMapper = map[string]TaxIdentifier{
		"Ђ": TaxPPdvNumericalIdentifier,
		"Е": TaxOPdvNumericalIdentifier,
	}
)

func (tid TaxIdentifier) Tax() *Tax {
	switch tid {
	case TaxPPdvNumericalIdentifier:
		return &TaxPPdv
	case TaxOPdvNumericalIdentifier:
		return &TaxOPdv
	default:
		return nil
	}
}
