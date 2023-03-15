package receipt_fetcher

import (
	"strconv"
	"strings"
	"time"
)

// Receipt DTO
type Receipt struct {
	Store               Store                `json:"store"`
	Number              string               `json:"number"`
	Counter             string               `json:"counter"`
	MetaData            MetaData             `json:"meta_data"`
	Items               []ReceiptItem        `json:"receipt_items"`
	Taxes               []TaxItem            `json:"tax_items"`
	PaymentSummary      map[string]RsdAmount `json:"payment_summary"`
	TotalPurchaseAmount RsdAmount            `json:"total_purchase_amount"`
	TotalTaxAmount      RsdAmount            `json:"total_tax_amount"`
	Date                time.Time            `json:"date"`
	QrCod               string               `json:"qr_code"`
}

// Store DTO
type Store struct {
	Name         string `json:"name"`
	Tin          string `json:"tin"`
	LocationId   string `json:"location_id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
	City         string `json:"city"`
}

// Stores Meta Data about the receipt
type MetaData map[string]string

// ReceiptItem DTO
type ReceiptItem struct {
	Name         string    `json:"name"`
	Quantity     float64   `json:"quantity"`
	Unit         string    `json:"unit"`
	Tax          Tax       `json:"tax"`
	SingleAmount RsdAmount `json:"single_amount"`
	TotalAmount  RsdAmount `json:"total_amount"`
}

// Tax DTO
type Tax struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Rate       int    `json:"rate"`
}

// TaxItemDTO
type TaxItem struct {
	Tax    Tax       `json:"tax"`
	Amount RsdAmount `json:"amount"`
}

// Fraction of real number
type RsdAmount struct {
	Value    int `json:"value"`
	Fraction int `json:"fraction"`
}

func (rsdAmount RsdAmount) GetFloat() float64 {
	return float64(rsdAmount.Value) + float64(rsdAmount.Fraction)/100
}

func (rsdAmount RsdAmount) GetParas() int {
	return rsdAmount.Value*100 + rsdAmount.Fraction
}

func NewRsdAmountFromString(str string) (RsdAmount, error) {
	rsdAmountData := strings.Split(str, ",")
	value, _ := strconv.Atoi(strings.Replace(rsdAmountData[0], ".", "", 2))
	fraction, _ := strconv.Atoi(rsdAmountData[1])

	rsdAmount := RsdAmount{
		Value:    value,
		Fraction: fraction,
	}

	return rsdAmount, nil
}
