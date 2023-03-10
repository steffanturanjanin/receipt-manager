package receipt_fetcher

import (
	"strconv"
	"strings"
	"time"
)

// Receipt DTO
type Receipt struct {
	Store               Store
	Number              string
	Counter             string
	Meta                Meta
	Items               []ReceiptItem
	Taxes               []TaxItem
	PaymentSummary      map[string]RsdAmount
	TotalPurchaseAmount RsdAmount
	TotalTaxAmount      RsdAmount
	Date                time.Time
	QrCod               string
}

// Store DTO
type Store struct {
	Name         string
	Tin          string
	LocationId   string
	LocationName string
	Address      string
	City         string
}

// Stores Meta Data about the receipt
type Meta map[string]string

// ReceiptItem DTO
type ReceiptItem struct {
	Name         string
	Quantity     float64
	Unit         string
	Tax          Tax
	SingleAmount RsdAmount
	TotalAmount  RsdAmount
}

// Tax DTO
type Tax struct {
	Name       string
	Identifier string
	Rate       int
}

// TaxItemDTO
type TaxItem struct {
	Tax    Tax
	Amount RsdAmount
}

// Fraction of real number
type RsdAmount struct {
	Value    int
	Fraction int
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
