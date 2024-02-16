package transport

import (
	"encoding/json"
	"time"

	"github.com/steffanturanjanin/receipt-manager/internal/models"
)

type ReceiptResponse struct {
	ID                  int                   `json:"id"`
	UserId              *int                  `json:"userId"`
	Store               *StoreResponse        `json:"store"`
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

type ReceiptsResponse struct {
	Items []ReceiptResponse
}

func (receipts ReceiptsResponse) FromModels(models []models.Receipt) ReceiptsResponse {
	for _, dbReceipt := range models {
		receipt := ReceiptResponseFromReceipt(dbReceipt)
		receipts.Items = append(receipts.Items, receipt)
	}

	return receipts
}

func (receipt ReceiptResponse) FromModel(model models.Receipt) ReceiptResponse {
	//userId := new(int)
	var userId *int
	if model.UserID != nil {
		userId = new(int)
		*userId = int(*model.UserID)
	}

	var store *StoreResponse
	if model.Store != nil {
		store = new(StoreResponse)
		*store = store.FromModel(*model.Store)
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
			receiptItem := ReceiptItemResponse{}
			receiptItem = receiptItem.FromModel(dbReceiptItem)
			receiptItems = append(receiptItems, receiptItem)
		}
	}

	receipt.ID = int(model.ID)
	receipt.UserId = userId
	receipt.Store = store
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

func ReceiptsResponseFromReceipt(models []models.Receipt) []ReceiptResponse {
	receipts := []ReceiptResponse{}

	for _, dbReceipt := range models {
		receipt := ReceiptResponseFromReceipt(dbReceipt)
		receipts = append(receipts, receipt)
	}

	return receipts
}

func ReceiptResponseFromReceipt(model models.Receipt) ReceiptResponse {
	receipt := ReceiptResponse{}

	userId := new(int)
	if model.UserID != nil {
		*userId = int(*model.UserID)
	}

	var store *StoreResponse
	if model.Store != nil {
		store.ID = int(model.Store.ID)
		store.Tin = model.Store.Tin
		store.LocationId = model.Store.LocationId
		store.LocationName = model.Store.LocationName
		store.Address = model.Store.Address
		store.City = model.Store.City
	}

	var meta *map[string]string
	if model.Meta != nil {
		if err := json.Unmarshal(*model.Meta, meta); err != nil {
			meta = nil
		}
	}

	receiptItems := make([]ReceiptItemResponse, 0)
	if model.ReceiptItems != nil {
		for _, dbReceiptItem := range model.ReceiptItems {
			var category *CategoryResponse
			if dbReceiptItem.Category != nil {
				category.ID = int(dbReceiptItem.Category.ID)
				category.Name = dbReceiptItem.Category.Name
			}

			var tax *TaxResponse
			if dbReceiptItem.Tax != nil {
				tax.ID = int(dbReceiptItem.Tax.ID)
				tax.Identifier = dbReceiptItem.Tax.Identifier
				tax.Name = dbReceiptItem.Tax.Name
				tax.Rate = dbReceiptItem.Tax.Rate
			}

			receiptItem := ReceiptItemResponse{}
			receiptItem.ID = int(dbReceiptItem.ID)
			receiptItem.ReceiptId = int(dbReceiptItem.ReceiptID)
			receiptItem.Name = dbReceiptItem.Name
			receiptItem.Unit = dbReceiptItem.Unit
			receiptItem.Quantity = dbReceiptItem.Quantity
			receiptItem.SingleAmount = dbReceiptItem.SingleAmount
			receiptItem.TotalAmount = dbReceiptItem.TotalAmount
			receiptItem.Category = category
			receiptItem.Tax = tax

			receiptItems = append(receiptItems, receiptItem)
		}
	}

	receipt.ID = int(model.ID)
	receipt.UserId = userId
	receipt.Store = store
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
