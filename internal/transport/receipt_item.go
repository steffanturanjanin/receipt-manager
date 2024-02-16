package transport

import "github.com/steffanturanjanin/receipt-manager/internal/models"

type ReceiptItemResponse struct {
	ID           int               `json:"id"`
	ReceiptId    int               `json:"receiptId"`
	Name         string            `json:"name"`
	Unit         string            `json:"unit"`
	Quantity     float64           `json:"quantity"`
	SingleAmount int               `json:"singleAmount"`
	TotalAmount  int               `json:"totalAmount"`
	Category     *CategoryResponse `json:"category"`
	Tax          *TaxResponse      `json:"tax"`
}

func (receiptItem ReceiptItemResponse) FromModel(model models.ReceiptItem) ReceiptItemResponse {
	var category *CategoryResponse
	if model.Category != nil {
		category = new(CategoryResponse)
		*category = category.FromModel(*model.Category)
	}
	var tax *TaxResponse
	if model.Tax != nil {
		tax = new(TaxResponse)
		*tax = tax.FromModel(*model.Tax)
	}

	receiptItem.ID = int(model.ID)
	receiptItem.ReceiptId = int(model.ReceiptID)
	receiptItem.Name = model.Name
	receiptItem.Unit = model.Unit
	receiptItem.Quantity = model.Quantity
	receiptItem.SingleAmount = model.SingleAmount
	receiptItem.TotalAmount = model.TotalAmount
	receiptItem.Category = category
	receiptItem.Tax = tax

	return receiptItem
}
