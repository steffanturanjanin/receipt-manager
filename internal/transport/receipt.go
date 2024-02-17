package transport

import (
	"encoding/json"
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/models"
)

type BaseReceiptResponse struct {
	ID                  int                   `json:"id"`
	UserId              *int                  `json:"userId"`
	Status              string                `json:"status"`
	PfrNumber           *string               `json:"pfrNumber"`
	Counter             *string               `json:"counter"`
	TotalPurchaseAmount *int                  `json:"totalPurchaseAmount"`
	TotalTaxAmount      *int                  `json:"totalTaxAmount"`
	Date                *time.Time            `json:"date"`
	QrCode              *string               `json:"qrCode"`
	Meta                *map[string]string    `json:"meta"`
	ReceiptItems        []ReceiptItemResponse `json:"receiptItems"`
}

type BaseReceiptTransformer struct{}

func (t BaseReceiptTransformer) Transform(models []models.Receipt) []BaseReceiptResponse {
	receipts := make([]BaseReceiptResponse, 0)

	for _, model := range models {
		receipt := t.TransformSingle(model)
		receipts = append(receipts, receipt)
	}

	return receipts
}

func (t BaseReceiptTransformer) TransformSingle(model models.Receipt) BaseReceiptResponse {
	receipt := BaseReceiptResponse{}

	var userId *int
	if model.UserID != nil {
		userId = new(int)
		*userId = int(*model.UserID)
	}

	var meta *map[string]string
	if model.Meta != nil {
		meta = new(map[string]string)
		if err := json.Unmarshal(*model.Meta, meta); err != nil {
			meta = nil
		}
	}

	receiptItems := make([]ReceiptItemResponse, 0)
	if model.ReceiptItems != nil {
		for _, dbReceiptItem := range model.ReceiptItems {
			transformer := ReceiptItemTransformer{}
			receiptItem := transformer.TransformSingle(dbReceiptItem)
			receiptItems = append(receiptItems, receiptItem)
		}
	}

	receipt.ID = int(model.ID)
	receipt.UserId = userId
	receipt.Status = model.Status
	receipt.PfrNumber = model.PfrNumber
	receipt.Counter = model.Counter
	receipt.TotalPurchaseAmount = model.TotalPurchaseAmount
	receipt.TotalTaxAmount = model.TotalTaxAmount
	receipt.Date = model.Date
	receipt.QrCode = model.QrCode
	receipt.Meta = meta
	receipt.ReceiptItems = receiptItems

	return receipt
}

type ReceiptResponse struct {
	BaseReceiptResponse
	Store *StoreResponse `json:"store"`
}

type ReceiptTransformer struct {
	BaseReceiptTransformer
}

func (t ReceiptTransformer) Transform(models []models.Receipt) []ReceiptResponse {
	items := make([]ReceiptResponse, 0)

	for _, model := range models {
		item := t.TransformSingle(model)
		items = append(items, item)
	}

	return items
}

func (t ReceiptTransformer) TransformSingle(model models.Receipt) ReceiptResponse {
	baseReceipt := t.BaseReceiptTransformer.TransformSingle(model)
	receipt := ReceiptResponse{BaseReceiptResponse: baseReceipt}

	var store *StoreResponse
	if model.Store != nil {
		transformer := StoreTransformer{}
		store = new(StoreResponse)
		*store = transformer.TransformSingle(*model.Store)
	}

	receipt.Store = store

	return receipt
}
